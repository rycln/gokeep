// Package password provides secure password hashing and verification.
package password

import (
	"golang.org/x/crypto/bcrypt"
)

// BCryptHasher implements password hashing using bcrypt algorithm.
type BCryptHasher struct{}

// NewBCryptHasher creates a new bcrypt password hasher instance.
func NewBCryptHasher() *BCryptHasher {
	return &BCryptHasher{}
}

// Hash generates a secure bcrypt hash from a plaintext password.
func (h *BCryptHasher) Hash(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedBytes), err
}

// Compare verifies a plaintext password against a stored hash.
func (h *BCryptHasher) Compare(hashed, plain string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	if err != nil {
		return newErrWrongPassword(ErrWrongPassword)
	}
	return nil
}
