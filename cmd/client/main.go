package main

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rycln/gokeep/internal/client/tui"
)

func main() {
	p := tea.NewProgram(tui.NewAuthModel())
	if _, err := p.Run(); err != nil {
		os.Exit(1)
	}
}
