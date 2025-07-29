package styles

import "github.com/charmbracelet/lipgloss"

var (
	ErrorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	TitleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("63"))
	InputStyle   = lipgloss.NewStyle().PaddingLeft(1)
	FocusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	ButtonStyle  = lipgloss.NewStyle().Padding(0, 3).Background(lipgloss.Color("62")).Foreground(lipgloss.Color("230"))
	ActiveButton = lipgloss.NewStyle().Background(lipgloss.Color("69"))
)
