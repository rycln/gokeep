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
	UserRegister(context.Context, *models.UserRegReq) (*models.User, error)
	UserLogin(context.Context, *models.UserLoginReq) (*models.User, error)
}

// saltGenerator defines operations for generating cryptographic salt
type saltGenerator interface {
	// GenerateSalt creates a new random salt value
	GenerateSalt() ([]byte, error)
}

// saltConverter defines operations for encoding and decoding salt values
type saltConverter interface {
	// EncodeSalt converts binary salt to string representation
	EncodeSalt([]byte) string
	// DecodeSalt converts string salt back to binary representation
	DecodeSalt(string) ([]byte, error)
}

// keyDeriver defines operations for deriving cryptographic keys
type keyDeriver interface {
	// DeriveKeyFromPasswordAndSalt creates a cryptographic key from password and salt
	DeriveKeyFromPasswordAndSalt(string, []byte) []byte
}

// keyProvider defines key handling for crypto operations, combining salt generation,
// conversion and key derivation capabilities
type keyProvider interface {
	saltGenerator
	saltConverter
	keyDeriver
}

// crypter defines interface for encryption and decryption operations
type crypter interface {
	// SetKey configures the encryption key to be used for subsequent operations
	SetKey(key []byte) error
}

// Model represents authentication screen state and its dependencies
type Model struct {
	state       state         // Current screen state
	activeField field         // Currently focused input field
	username    string        // Username input value
	password    string        // Password input value
	errMsg      string        // Last error message to display
	service     authService   // Authentication service implementation
	key         keyProvider   // Key generation and handling provider
	crypt       crypter       // Cryptographic operations handler
	timeout     time.Duration // Maximum duration for authentication operations
}

// InitialModel creates new authentication model with dependencies
func InitialModel(
	service authService,
	key keyProvider,
	crypt crypter,
	timeout time.Duration,
) Model {
	return Model{
		state:       LoginState,
		activeField: UsernameField,
		service:     service,
		timeout:     timeout,
		key:         key,
		crypt:       crypt,
	}
}

// GetState returns current authentication screen state
// Returns the current state enum value
func (m Model) GetState() state {
	return m.state
}
