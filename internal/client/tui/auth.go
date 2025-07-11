package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type AuthModel struct {
	usernameInput textinput.Model // Поле для логина
	passwordInput textinput.Model // Поле для пароля
	focusedField  int             // 0 = логин, 1 = пароль
	err           error
}

func (m AuthModel) Init() tea.Cmd {
	return nil
}

func NewAuthModel() AuthModel {
	username := textinput.New()
	username.Placeholder = "Логин"
	username.Focus()

	password := textinput.New()
	password.Placeholder = "Пароль"
	password.EchoMode = textinput.EchoPassword

	return AuthModel{
		usernameInput: username,
		passwordInput: password,
		focusedField:  0,
	}
}

func (m AuthModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab", "shift+tab", "enter":
			/*
				if msg.String() == "enter" && m.focusedField == 1 {
					// Вызов логики авторизации
					err := app.Login(m.usernameInput.Value(), m.passwordInput.Value())
					if err != nil {
						m.err = err
						return m, nil
					}
					return NewListModel(), nil // Переход к следующему экрану
				}
			*/

			if msg.String() == "tab" || (msg.String() == "enter" && m.focusedField == 0) {
				m.focusedField = (m.focusedField + 1) % 2
			} else {
				m.focusedField = (m.focusedField - 1) % 2
			}

			if m.focusedField == 0 {
				m.usernameInput.Focus()
				m.passwordInput.Blur()
			} else {
				m.usernameInput.Blur()
				m.passwordInput.Focus()
			}

		default:
			var cmd tea.Cmd
			if m.focusedField == 0 {
				m.usernameInput, cmd = m.usernameInput.Update(msg)
			} else {
				m.passwordInput, cmd = m.passwordInput.Update(msg)
			}
			return m, cmd
		}
	}
	return m, nil
}

func (m AuthModel) View() string {
	return fmt.Sprintf(
		"Вход в GophKeeper\n\n%s\n%s\n\n%s",
		m.usernameInput.View(),
		m.passwordInput.View(),
		"(tab - переключение, enter - подтвердить)",
	) //+ "\n\n" + m.err.Error()
}
