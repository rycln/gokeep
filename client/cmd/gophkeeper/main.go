package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"
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

const (
	DBpath  = "./gophkeeper.db"
	timeout = time.Duration(5) * time.Second
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

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err = storage.InitDB(ctx, db)
	if err != nil {
		log.Fatal(err)
	}

	itemStorage := storage.NewItemStorage(db)

	authService := services.NewAuthService(client.NewUserClient(conn))
	itemService := services.NewItemService(itemStorage)

	authScreen := auth.InitialModel(authService, timeout)
	vaultScreen := vault.InitialModel(itemService, timeout)
	addScreen := add.InitialModel(itemService, timeout)
	updateScreen := update.InitialModel(itemService, timeout)

	p := tea.NewProgram(tui.InitialRootModel(authScreen, vaultScreen, addScreen, updateScreen))
	if _, err := p.Run(); err != nil {
		os.Exit(1)
	}
}
