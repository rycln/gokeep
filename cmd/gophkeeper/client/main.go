package main

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	client "github.com/rycln/gokeep/internal/client/grpc"
	"github.com/rycln/gokeep/internal/client/services"
	"github.com/rycln/gokeep/internal/client/tui"
	"github.com/rycln/gokeep/internal/client/tui/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

	authService := services.NewAuthService(client.NewAuthClient(conn))

	authScreen := auth.InitialModel(authService)

	p := tea.NewProgram(tui.InitialRootModel(authScreen))
	if _, err := p.Run(); err != nil {
		os.Exit(1)
	}
}
