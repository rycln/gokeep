package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	client "github.com/rycln/gokeep/internal/client/grpc"
	"github.com/rycln/gokeep/internal/client/services"
	"github.com/rycln/gokeep/internal/client/storage"
	"github.com/rycln/gokeep/internal/client/tui"
	"github.com/rycln/gokeep/internal/client/tui/auth"
	"github.com/rycln/gokeep/internal/client/tui/vault"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	DBpath = "./gophkeeper.db"
)

func main() {
	//временно в main
	certPool, _ := x509.SystemCertPool()

	tlsConfig := &tls.Config{
		RootCAs:    certPool,
		MinVersion: tls.VersionTLS12,
	}
	conn, err := grpc.NewClient(":50051", grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	db, err := storage.NewDB(DBpath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	itemStorage := storage.NewItemStorage(db)

	authService := services.NewAuthService(client.NewUserClient(conn))
	itemService := services.NewItemService(itemStorage)

	authScreen := auth.InitialModel(authService)
	vaultScreen := vault.NewModel(itemService, context.Background())

	p := tea.NewProgram(tui.InitialRootModel(authScreen, vaultScreen))
	if _, err := p.Run(); err != nil {
		os.Exit(1)
	}
}
