// Package crypto provides AES-GCM encryption/decryption functionality.
// Implements secure symmetric encryption for sensitive data storage.
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

// Error definitions for crypto operations
var (
	errShortKeySize = errors.New("invalid key size, expected 32 bytes") // AES-256 requires 32-byte key
	errShortCrypted = errors.New("crypted content too short")           // Minimum valid ciphertext length
	errNoKey        = errors.New("no key")                              // Key not initialized
)

// AESCrypter implements AES-GCM encryption/decryption
// Uses 256-bit keys and provides authenticated encryption
type AESCrypter struct {
	key []byte // Encryption key (must be 32 bytes)
}

// NewAESCrypter creates a new AES crypter instance
// Note: Key must be set via SetKey before use
func NewAESCrypter() *AESCrypter {
	return &AESCrypter{}
}

// SetKey configures the encryption key
func (c *AESCrypter) SetKey(key []byte) error {
	if len(key) != 32 {
		return errShortKeySize
	}
	c.key = key
	return nil
}

// Encrypt encrypts data using AES-GCM
func (c *AESCrypter) Encrypt(content []byte) ([]byte, error) {
	if len(c.key) == 0 {
		return nil, errNoKey
	}

	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, content, nil), nil
}

// Decrypt decrypts AES-GCM encrypted data
func (c *AESCrypter) Decrypt(crypted []byte) ([]byte, error) {
	if len(c.key) == 0 {
		return nil, errNoKey
	}

	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(crypted) < nonceSize {
		return nil, errShortCrypted
	}

	nonce, ciphertext := crypted[:nonceSize], crypted[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
