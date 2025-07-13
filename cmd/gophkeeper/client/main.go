package main

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rycln/gokeep/internal/client/grpc"
	"github.com/rycln/gokeep/internal/client/services"
	"github.com/rycln/gokeep/internal/client/tui"
	"github.com/rycln/gokeep/internal/client/tui/auth"
)

func main() {
	authService := services.NewAuthService(grpc.NewAuthClient())

	authScreen := auth.InitialModel(authService)

	p := tea.NewProgram(tui.InitialRootModel(authScreen))
	if _, err := p.Run(); err != nil {
		os.Exit(1)
	}
}
