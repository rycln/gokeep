package auth

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rycln/gokeep/shared/models"
)

// Init initializes the authentication model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles all messages and state transitions for authentication
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case LoginState, RegisterState:
		return handleAuthInput(m, msg)
	case ProcessingState:
		return handleProcessingState(m, msg)
	case ErrorState:
		return handleErrorState(m, msg)
	}
	return m, nil
}

// handleAuthInput processes user input in login/register states
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
			if m.activeField == UsernameField {
				m.activeField = PasswordField
			} else {
				m.activeField = UsernameField
			}
		case tea.KeyRunes:
			if msg.String() == " " {
				return m, nil
			}
			if m.activeField == UsernameField {
				m.username += msg.String()
			} else {
				m.password += msg.String()
			}
		case tea.KeyBackspace:
			if m.activeField == UsernameField && len(m.username) > 0 {
				runes := []rune(m.username)
				if len(runes) > 0 {
					m.username = string(runes[:len(runes)-1])
				}
			} else if len(m.password) > 0 {
				runes := []rune(m.password)
				if len(runes) > 0 {
					m.password = string(runes[:len(runes)-1])
				}
			}
		}
	}
	return m, nil
}

// login initiates user authentication
func (m Model) login() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
		defer cancel()

		user, err := m.service.UserLogin(ctx, &models.UserLoginReq{
			Username: m.username,
			Password: m.password,
		})
		if err != nil {
			return LoginErrorMsg{err}
		}

		decSalt, err := m.key.DecodeSalt(user.Salt)
		if err != nil {
			return LoginErrorMsg{err}
		}

		key := m.key.DeriveKeyFromPasswordAndSalt(m.password, decSalt)

		err = m.crypt.SetKey(key)
		if err != nil {
			return LoginErrorMsg{err}
		}

		return AuthSuccessMsg{user}
	}
}

// register initiates new user registration
func (m Model) register() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
		defer cancel()

		salt, err := m.key.GenerateSalt()
		if err != nil {
			return RegisterErrorMsg{err}
		}

		encSalt := m.key.EncodeSalt(salt)

		user, err := m.service.UserRegister(ctx, &models.UserRegReq{
			Username: m.username,
			Password: m.password,
			Salt:     encSalt,
		})
		if err != nil {
			return RegisterErrorMsg{err}
		}

		key := m.key.DeriveKeyFromPasswordAndSalt(m.password, salt)

		err = m.crypt.SetKey(key)
		if err != nil {
			return RegisterErrorMsg{err}
		}

		return AuthSuccessMsg{user}
	}
}

// handleProcessingState manages authentication operation results
func handleProcessingState(m Model, msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case LoginErrorMsg:
		m.errMsg = msg.Err.Error()
		m.state = ErrorState
	case RegisterErrorMsg:
		m.errMsg = msg.Err.Error()
		m.state = ErrorState
	case AuthSuccessMsg:
		return m, func() tea.Msg { return msg }
	}
	return m, nil
}

// handleErrorState manages error display and recovery
func handleErrorState(m Model, msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.state = LoginState
		}
	}
	return m, nil
}
