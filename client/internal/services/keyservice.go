package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"

	"golang.org/x/crypto/pbkdf2"
)

// Security parameters for key derivation
const (
	keyLength        = 32     // 256-bit key length for AES-256
	saltLength       = 16     // 128-bit salt length
	pbkdf2Iterations = 600000 // NIST recommended minimum iterations
)

var errInvalidSaltLength = errors.New("invalid salt length")

// KeyService handles cryptographic key operations
type KeyService struct{}

// NewKeyService creates a new KeyService instance
func NewKeyService() *KeyService {
	return &KeyService{}
}

// GenerateSalt creates a new cryptographically secure random salt
func (s *KeyService) GenerateSalt() ([]byte, error) {
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	return salt, nil
}

// DeriveKeyFromPasswordAndSalt generates a cryptographic key using PBKDF2
func (s *KeyService) DeriveKeyFromPasswordAndSalt(password string, salt []byte) []byte {
	return pbkdf2.Key(
		[]byte(password),
		salt,
		pbkdf2Iterations,
		keyLength,
		sha256.New,
	)
}

// DecodeSalt converts base64-encoded salt back to bytes
func (s *KeyService) DecodeSalt(salt string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return nil, err
	}
	if len(decoded) != saltLength {
		return nil, errInvalidSaltLength
	}
	return decoded, nil
}

// EncodeSalt converts salt bytes to base64 string
func (s *KeyService) EncodeSalt(salt []byte) string {
	return base64.StdEncoding.EncodeToString(salt)
}
