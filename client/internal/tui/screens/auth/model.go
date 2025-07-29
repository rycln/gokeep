package auth

import (
	"context"
	"time"

	"github.com/rycln/gokeep/shared/models"
)

type state int

const (
	LoginState state = iota
	RegisterState
	ProcessingState
	ErrorState
)

type field int

const (
	UsernameField field = iota
	PasswordField
)

type (
	AuthSuccessMsg   struct{ User *models.User }
	LoginErrorMsg    struct{ Err error }
	RegisterErrorMsg struct{ Err error }
)

type authService interface {
	UserRegister(context.Context, *models.UserAuthReq) (*models.User, error)
	UserLogin(context.Context, *models.UserAuthReq) (*models.User, error)
}

type Model struct {
	state       state
	activeField field
	username    string
	password    string
	errMsg      string
	service     authService
	timeout     time.Duration
}

func InitialModel(service authService, timeout time.Duration) Model {
	return Model{
		state:       LoginState,
		activeField: UsernameField,
		service:     service,
		timeout:     timeout,
	}
}

func (m Model) GetState() state {
	return m.state
}
