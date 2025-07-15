package main

import (
	"log"
	"net"
	"time"

	"github.com/rycln/gokeep/internal/server/grpc/server"
	"github.com/rycln/gokeep/internal/server/services"
	"github.com/rycln/gokeep/internal/server/storage"
	"github.com/rycln/gokeep/internal/server/strategies/password"
	pb "github.com/rycln/gokeep/internal/shared/proto/gen/gophkeeper"
	"google.golang.org/grpc"
)

const (
	basePort = ":50051"
	dsn      = "some dsn"
	jwtKey   = "secret"
)

var jwtExp = time.Duration(1) * time.Hour

func main() {
	db, err := storage.NewDB(dsn)
	if err != nil {
		log.Fatal(err)
	}

	authstrg := storage.NewUserStorage(db)

	passwordStrategy := password.NewBCryptHasher()
	jwtservice := services.NewJWTService(jwtKey, jwtExp)
	authservice := services.NewUserService(authstrg, passwordStrategy, jwtservice)

	g := grpc.NewServer()

	gs := server.NewGophKeeperServer(authservice)

	pb.RegisterGophKeeperServer(g, gs)

	listen, err := net.Listen("tcp", basePort)
	if err != nil {
		log.Fatalf("gRPC server error: %v", err)
	}
	if err := g.Serve(listen); err != nil {
		log.Fatalf("gRPC server error: %v", err)
	}
}
