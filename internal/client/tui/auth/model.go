package auth

import (
	"context"

	"github.com/rycln/gokeep/internal/shared/models"
)

type State int

const (
	LoginState State = iota
	RegisterState
	ProcessingState
	ErrorState
)

type authService interface {
	UserRegister(context.Context, *models.UserAuthReq) (*models.User, error)
	UserLogin(context.Context, *models.UserAuthReq) (*models.User, error)
}

type Model struct {
	state       State
	activeField string
	username    string
	password    string
	errMsg      string
	service     authService
}

type (
	AuthSuccessMsg   struct{ User *models.User }
	LoginErrorMsg    struct{ Err error }
	RegisterErrorMsg struct{ Err error }
)

func InitialModel(service authService) Model {
	return Model{
		state:       LoginState,
		activeField: "username",
		service:     service,
	}
}

func (m Model) GetState() State {
	return m.state
}
