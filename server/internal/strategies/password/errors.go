package password

import "errors"

// ErrWrongPassword indicates a password verification failure.
var ErrWrongPassword = errors.New("wrong password")

// errWrongPassword implements a structured password verification error.
type errWrongPassword struct {
	err error // Underlying error
}

// Error implements the error interface.
func (err *errWrongPassword) Error() string {
	return err.err.Error()
}

// Unwrap supports error inspection with errors.Is()/errors.As().
func (err *errWrongPassword) Unwrap() error {
	return err.err
}

// IsErrWrongPassword provides type checking method.
func (err *errWrongPassword) IsErrWrongPassword() bool {
	return true
}

// newErrWrongPassword constructs a new password verification error.
func newErrWrongPassword(err error) error {
	return &errWrongPassword{
		err: err,
	}
}
