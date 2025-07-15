package auth

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rycln/gokeep/internal/shared/models"
)

const timeout = time.Duration(5) * time.Second

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case LoginState, RegisterState:
		return handleAuthInput(m, msg)
	case ProcessingState:
		switch msg := msg.(type) {
		case LoginErrorMsg:
			m.errMsg = msg.Err.Error()
			m.state = ErrorState
		case RegisterErrorMsg:
			m.errMsg = msg.Err.Error()
			m.state = ErrorState
		case AuthSuccessMsg:
			m.state = SuccessState
		}
		return m, nil
	case ErrorState:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				m.state = LoginState
			}
		}
		return m, nil
	case SuccessState:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				m.state = LoginState
			}
		}
	}
	return m, nil
}

func handleAuthInput(m Model, msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			if m.state == LoginState {
				m.state = ProcessingState
				return m, m.login()
			} else {
				m.state = ProcessingState
				return m, m.register()
			}
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
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		user, err := m.service.UserLogin(ctx, &models.UserAuthReq{
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
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		user, err := m.service.UserRegister(ctx, &models.UserAuthReq{
			Username: m.username,
			Password: m.password,
		})
		if err != nil {
			return RegisterErrorMsg{err}
		}
		return AuthSuccessMsg{user}
	}
}
