package auth

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rycln/gokeep/internal/shared/models"
)

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case LoginState, RegisterState:
		return handleAuthInput(m, msg)
	case ProcessingState:
		return m, nil
	default:
		//здесь сделать успешный вход
		return m, nil
	}
}

func handleAuthInput(m Model, msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			m.state = ProcessingState
			if m.state == LoginState {
				return m, m.login()
			}
			return m, m.register()
		case tea.KeyTab:
			if m.state == LoginState {
				m.state = RegisterState
			} else {
				m.state = LoginState
			}
			return m, nil
		case tea.KeyDown, tea.KeyUp:
			if m.activeField == "username" {
				m.activeField = "password"
			} else {
				m.activeField = "username"
			}
		case tea.KeyRunes:
			if msg.String() == " " {
				return m, nil
			}
			if m.activeField == "username" {
				m.username += msg.String()
			} else {
				m.password += msg.String()
			}
		case tea.KeyBackspace:
			if m.activeField == "username" && len(m.username) > 0 {
				m.username = m.username[:len(m.username)-1]
			} else if len(m.password) > 0 {
				m.password = m.password[:len(m.password)-1]
			}
		}
	}
	return m, nil
}

func (m Model) login() tea.Cmd {
	return func() tea.Msg {
		user, err := m.service.UserLogin(context.Background(), &models.UserAuthReq{
			Username: m.username,
			Password: m.password,
		})
		if err != nil {
			return LoginErrorMsg{err}
		}
		return AuthSuccessMsg{user}
	}
}

func (m Model) register() tea.Cmd {
	return func() tea.Msg {
		user, err := m.service.UserRegister(context.Background(), &models.UserAuthReq{
			Username: m.username,
			Password: m.password,
		})
		if err != nil {
			return RegisterErrorMsg{err}
		}
		return AuthSuccessMsg{user}
	}
}
