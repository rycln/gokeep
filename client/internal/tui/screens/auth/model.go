// Package auth implements authentication screen logic.
package auth

import (
	"context"
	"time"

	"github.com/rycln/gokeep/shared/models"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

// state represents current authentication screen state
type state int

// Authentication screen states
const (
	LoginState      state = iota // User login form
	RegisterState                // User registration form
	ProcessingState              // Authentication in progress
	ErrorState                   // Error display state
)

// field represents active form field
type field int

// Form field constants
const (
	UsernameField field = iota // Username input field
	PasswordField              // Password input field
)

// Message types for authentication events
type (
	// AuthSuccessMsg indicates successful authentication
	AuthSuccessMsg struct{ User *models.User }

	// LoginErrorMsg contains login failure details
	LoginErrorMsg struct{ Err error }

	// RegisterErrorMsg contains registration failure details
	RegisterErrorMsg struct{ Err error }
)

// authService defines required authentication operations
type authService interface {
	UserRegister(context.Context, *models.UserAuthReq) (*models.User, error)
	UserLogin(context.Context, *models.UserAuthReq) (*models.User, error)
}

// Model represents authentication screen state
type Model struct {
	state       state         // Current screen state
	activeField field         // Currently focused field
	username    string        // Username input value
	password    string        // Password input value
	errMsg      string        // Last error message
	service     authService   // Authentication service
	timeout     time.Duration // Operation timeout
}

// InitialModel creates new authentication model with dependencies
func InitialModel(service authService, timeout time.Duration) Model {
	return Model{
		state:       LoginState,
		activeField: UsernameField,
		service:     service,
		timeout:     timeout,
	}
}

// GetState returns current authentication screen state
func (m Model) GetState() state {
	return m.state
}
