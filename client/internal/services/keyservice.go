package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"

	"golang.org/x/crypto/pbkdf2"
)

const (
	keyLength        = 32
	saltLength       = 16
	pbkdf2Iterations = 600000
)

var errInvalidSaltLength = errors.New("invalid salt length")

type KeyService struct{}

func NewKeyService() *KeyService {
	return &KeyService{}
}

func (s *KeyService) GenerateSalt() ([]byte, error) {
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	return salt, nil
}

func (s *KeyService) DeriveKeyFromPassword(password string, salt []byte) ([]byte, error) {
	return pbkdf2.Key(
		[]byte(password),
		salt,
		pbkdf2Iterations,
		keyLength,
		sha256.New,
	), nil
}

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

func (s *KeyService) EncodeSalt(salt []byte) string {
	return base64.StdEncoding.EncodeToString(salt)
}
