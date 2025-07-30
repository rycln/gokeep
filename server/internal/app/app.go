// Package app is the root package that composes all application components
// into a runnable service.
package app

import (
	"crypto/tls"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	pb "github.com/rycln/gokeep/pkg/gen/grpc/gophkeeper"
	"github.com/rycln/gokeep/server/internal/config"
	server "github.com/rycln/gokeep/server/internal/grpc"
	"github.com/rycln/gokeep/server/internal/grpc/interceptors"
	"github.com/rycln/gokeep/server/internal/logger"
	"github.com/rycln/gokeep/server/internal/services"
	"github.com/rycln/gokeep/server/internal/storage"
	"github.com/rycln/gokeep/server/internal/strategies/password"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// buildInfo holds application build metadata that can be set during compilation.
var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

// Package-level constants defining core application parameters.
const (
	// jwtExpires sets the lifetime duration for JWT authentication tokens.
	// Used in auth service when generating new tokens.
	jwtExpires = time.Duration(2) * time.Hour
)

// App represents the core application layer.
type App struct {
	grpcserver *grpc.Server
	db         *sql.DB
	cfg        *config.Cfg
}

// New constructs and initializes the complete application.
func New() (*App, error) {
	cfg, err := config.NewConfigBuilder().
		WithFlagParsing().
		WithEnvParsing().
		WithConfigFile().
		WithDefaultJWTKey().
		Build()
	if err != nil {
		return nil, fmt.Errorf("can't initialize config: %v", err)
	}

	err = logger.LogInit(cfg.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("can't initialize logger: %v", err)
	}

	db, err := storage.NewDB(cfg.DatabaseDsn)
	if err != nil {
		return nil, fmt.Errorf("can't init DB: %v", err)
	}

	authstrg := storage.NewUserStorage(db)

	passwordStrategy := password.NewBCryptHasher()
	jwtservice := services.NewJWTService(cfg.Key, jwtExpires)
	authservice := services.NewUserService(authstrg, passwordStrategy, jwtservice)

	serverCert, err := tls.LoadX509KeyPair(cfg.CertFileName, cfg.CertKeyFileName)
	if err != nil {
		return nil, fmt.Errorf("can't load cert: %v", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
		MinVersion:   tls.VersionTLS12,
	}

	g := grpc.NewServer(
		grpc.Creds(credentials.NewTLS(tlsConfig)),
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(interceptors.InterceptorLogger(logger.Log)),
			auth.UnaryServerInterceptor(interceptors.NewAuthInterceptor(jwtservice).AuthFunc),
		),
	)

	gs := server.NewGophKeeperServer(authservice, cfg.Timeout)

	pb.RegisterUserServiceServer(g, gs)

	return &App{
		grpcserver: g,
		db:         db,
		cfg:        cfg,
	}, nil
}

// Run starts the application services.
func (app *App) Run() error {
	go func() {
		listen, err := net.Listen("tcp", app.cfg.GRPCPort)
		if err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
		if err := app.grpcserver.Serve(listen); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	logger.Log.Info(fmt.Sprintf("Server started successfully! Port: %s", app.cfg.GRPCPort))
	printBuildInfo()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	<-shutdown

	err := app.shutdown()
	if err != nil {
		return fmt.Errorf("shutdown error: %v", err)
	}

	err = app.cleanup()
	if err != nil {
		return fmt.Errorf("cleanup error: %v", err)
	}

	log.Println(strings.TrimPrefix(os.Args[0], "./") + " shutted down gracefully")

	return nil
}

// shutdown gracefully shuts down the application components.
func (app *App) shutdown() error {
	app.grpcserver.GracefulStop()

	return nil
}

// cleanup performs resource cleanup operations for the application.
func (app *App) cleanup() error {
	if err := app.db.Close(); err != nil {
		return fmt.Errorf("storage close failed: %w", err)
	}

	if err := logger.Log.Sync(); err != nil && !errors.Is(err, syscall.EINVAL) {
		return fmt.Errorf("log sync failed: %w", err)
	}

	return nil
}

// printBuildInfo displays the build metadata in a standardized format.
func printBuildInfo() {
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildDate == "" {
		buildDate = "N/A"
	}
	if buildCommit == "" {
		buildCommit = "N/A"
	}

	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}
