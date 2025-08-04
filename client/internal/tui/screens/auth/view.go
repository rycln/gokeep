package auth

import (
	"fmt"
	"strings"

	"github.com/rycln/gokeep/client/internal/tui/shared/i18n"
	"github.com/rycln/gokeep/client/internal/tui/shared/styles"
)

// View renders the current authentication screen based on state.
// Returns formatted UI with appropriate styling and localization.
func (m Model) View() string {
	switch m.state {
	case ProcessingState:
		return i18n.CommonWait
	case ErrorState:
		return styles.ErrorStyle.Render(fmt.Sprintf(i18n.CommonError, m.errMsg))
	default:
		return renderAuthForm(m)
	}
}

// renderAuthForm builds the authentication form UI
// Includes username/password fields and mode toggle buttons
func renderAuthForm(m Model) string {
	title := i18n.AuthRegisterTitle
	if m.state == LoginState {
		title = i18n.AuthLoginTitle
	}

	usernameInput := styles.InputStyle.Render(
		fmt.Sprintf(i18n.AuthUsernameLabel, m.username))
	passwordInput := styles.InputStyle.Render(
		fmt.Sprintf(i18n.AuthPasswordLabel, maskPassword(m.password)))

	if m.activeField == UsernameField {
		usernameInput = styles.FocusedStyle.Render(
			"> " + fmt.Sprintf(i18n.AuthUsernameLabel, m.username))
	} else {
		passwordInput = styles.FocusedStyle.Render(
			"> " + fmt.Sprintf(i18n.AuthPasswordLabel, maskPassword(m.password)))
	}

	loginBtn := styles.ButtonStyle.Render(i18n.AuthLoginButton)
	registerBtn := styles.ButtonStyle.Render(i18n.AuthRegisterButton)
	if m.state == LoginState {
		loginBtn = styles.ActiveButton.Render(i18n.AuthLoginButton)
	} else {
		registerBtn = styles.ActiveButton.Render(i18n.AuthRegisterButton)
	}

	return fmt.Sprintf(
		"%s\n\n%s\n%s\n\n%s %s\n\n%s",
		styles.TitleStyle.Render(title),
		usernameInput,
		passwordInput,
		loginBtn,
		registerBtn,
		i18n.AuthTabHint,
	)
}

// maskPassword obscures password input for display
func maskPassword(pwd string) string {
	return strings.Repeat("•", len(pwd))
}
