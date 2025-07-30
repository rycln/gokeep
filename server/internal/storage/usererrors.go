package storage

import "errors"

// Base error definitions for user-related operations
var (
	// ErrUsernameConflict indicates a duplicate username violation
	ErrUsernameConflict = errors.New("username already registered")

	// ErrNoUser indicates a missing user record
	ErrNoUser = errors.New("user does not exist")
)

// errUsernameConflict implements a structured conflict error
type errUsernameConflict struct {
	err error // Underlying error
}

// Error implements the error interface
func (err *errUsernameConflict) Error() string {
	return err.err.Error()
}

// Unwrap supports error inspection with errors.Is()/errors.As()
func (err *errUsernameConflict) Unwrap() error {
	return err.err
}

// IsErrUsernameConflict provides type checking method
func (err *errUsernameConflict) IsErrUsernameConflict() bool {
	return true
}

// newErrUsernameConflict constructs a new username conflict error
func newErrUsernameConflict(err error) error {
	return &errUsernameConflict{
		err: err,
	}
}

// errNoUser implements a structured "user not found" error
type errNoUser struct {
	err error // Underlying error
}

// Error implements the error interface
func (err *errNoUser) Error() string {
	return err.err.Error()
}

// Unwrap supports error inspection with errors.Is()/errors.As()
func (err *errNoUser) Unwrap() error {
	return err.err
}

// IsErrNoUser provides type checking method
func (err *errNoUser) IsErrNoUser() bool {
	return true
}

// newErrNoUser constructs a new user not found error
func newErrNoUser(err error) error {
	return &errNoUser{
		err: err,
	}
}
