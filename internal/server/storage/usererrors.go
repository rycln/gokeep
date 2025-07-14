package storage

import (
	"errors"
)

var (
	ErrUsernameConflict = errors.New("username already registered")
	ErrNoUser           = errors.New("user does not exist")
)

type errUsernameConflict struct {
	err error
}

func (err *errUsernameConflict) Error() string {
	return err.err.Error()
}

func (err *errUsernameConflict) Unwrap() error {
	return err.err
}

func (err *errUsernameConflict) IsErrUsernameConflict() bool {
	return true
}

func newErrUsernameConflict(err error) error {
	return &errUsernameConflict{
		err: err,
	}
}

type errNoUser struct {
	err error
}

func (err *errNoUser) Error() string {
	return err.err.Error()
}

func (err *errNoUser) Unwrap() error {
	return err.err
}

func (err *errNoUser) IsErrNoUser() bool {
	return true
}

func newErrNoUser(err error) error {
	return &errNoUser{
		err: err,
	}
}
