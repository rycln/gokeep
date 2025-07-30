// App package contains main application logic and initialization.
// Handles gRPC connections, database setup and TUI program lifecycle.
package app

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	client "github.com/rycln/gokeep/client/internal/grpc"
	"github.com/rycln/gokeep/client/internal/services"
	"github.com/rycln/gokeep/client/internal/storage"
	"github.com/rycln/gokeep/client/internal/tui"
	"github.com/rycln/gokeep/client/internal/tui/screens/add"
	"github.com/rycln/gokeep/client/internal/tui/screens/auth"
	"github.com/rycln/gokeep/client/internal/tui/screens/update"
	"github.com/rycln/gokeep/client/internal/tui/screens/vault"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Build info variables populated during compilation
var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

// Application constants
const (
	grpcTarget = ":50051"
	DBpath     = "./gophkeeper.db"
	timeout    = time.Duration(5) * time.Second
)

// App represents the main application structure
type App struct {
	tui  *tea.Program     // TUI program instance
	conn *grpc.ClientConn // gRPC connection
	db   *sql.DB          // Database connection
}

// New creates and initializes a new App instance
func New() (*App, error) {
	certPool, _ := x509.SystemCertPool()

	tlsConfig := &tls.Config{
		RootCAs:    certPool,
		MinVersion: tls.VersionTLS12,
	}
	conn, err := grpc.NewClient(grpcTarget, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	if err != nil {
		return nil, fmt.Errorf("grpc client error: %v", err)
	}

	db, err := storage.NewDB(DBpath)
	if err != nil {
		return nil, fmt.Errorf("db creation error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err = storage.InitDB(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("db init error: %v", err)
	}

	itemStorage := storage.NewItemStorage(db)

	authService := services.NewAuthService(client.NewUserClient(conn))
	itemService := services.NewItemService(itemStorage)

	authScreen := auth.InitialModel(authService, timeout)
	vaultScreen := vault.InitialModel(itemService, timeout)
	addScreen := add.InitialModel(itemService, timeout)
	updateScreen := update.InitialModel(itemService, timeout)

	p := tea.NewProgram(tui.InitialRootModel(authScreen, vaultScreen, addScreen, updateScreen))

	return &App{
		tui:  p,
		conn: conn,
		db:   db,
	}, nil
}

// Run starts the application and manages its lifecycle
func (app *App) Run() error {
	printBuildInfo()

	shutdown := make(chan struct{}, 1)

	go func() {
		if _, err := app.tui.Run(); err != nil {
			os.Exit(1)
		}
		close(shutdown)
	}()

	<-shutdown

	err := app.cleanup()
	if err != nil {
		return fmt.Errorf("cleanup error: %v", err)
	}

	log.Println(strings.TrimPrefix(os.Args[0], "./") + " shutted down gracefully")

	return nil
}

// cleanup handles resource cleanup on application shutdown
func (app *App) cleanup() error {
	if err := app.db.Close(); err != nil {
		return fmt.Errorf("storage close failed: %w", err)
	}

	if err := app.conn.Close(); err != nil {
		return fmt.Errorf("conn close failed: %w", err)
	}

	return nil
}

// printBuildInfo displays version information
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
