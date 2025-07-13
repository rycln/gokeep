package auth

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("63"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	inputStyle   = lipgloss.NewStyle().PaddingLeft(1)
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	buttonStyle  = lipgloss.NewStyle().Padding(0, 3).Background(lipgloss.Color("62")).Foreground(lipgloss.Color("230"))
	activeButton = lipgloss.NewStyle().Background(lipgloss.Color("69"))
)

func (m Model) View() string {
	switch m.state {
	case ProcessingState:
		return "Пожалуйста, подождите..."
	case SuccessState:
		return "✅ Успешно! Нажмите Enter для продолжения..."
	case ErrorState:
		return errorStyle.Render("Ошибка: " + m.errMsg)
	default:
		return renderAuthForm(m)
	}
}

func renderAuthForm(m Model) string {
	var title string
	if m.state == LoginState {
		title = "Вход в GophKeeper"
	} else {
		title = "Регистрация"
	}

	usernameInput := inputStyle.Render("Логин: " + m.username)
	passwordInput := inputStyle.Render("Пароль: " + maskPassword(m.password))

	if m.activeField == "username" {
		usernameInput = focusedStyle.Render("> Логин: " + m.username)
	} else {
		passwordInput = focusedStyle.Render("> Пароль: " + maskPassword(m.password))
	}

	loginBtn := buttonStyle.Render("Вход")
	registerBtn := buttonStyle.Render("Регистрация")
	if m.state == LoginState {
		loginBtn = activeButton.Render("Вход")
	} else {
		registerBtn = activeButton.Render("Регистрация")
	}

	return fmt.Sprintf(
		"%s\n\n%s\n%s\n\n%s %s\n\n%s",
		titleStyle.Render(title),
		usernameInput,
		passwordInput,
		loginBtn,
		registerBtn,
		"Нажмите Enter для подтверждения, Tab для переключения",
	)
}

func maskPassword(pwd string) string {
	return strings.Repeat("•", len(pwd))
}
