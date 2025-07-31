package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

var (
	errShortKeySize = errors.New("invalid key size, expected 32 bytes")
	errShortCrypted = errors.New("crypted content too short")
	errNoKey        = errors.New("no key")
)

type AESCrypter struct {
	key []byte
}

func NewAESCrypter() *AESCrypter { return &AESCrypter{} }

func (c *AESCrypter) SetKey(key []byte) error {
	if len(key) != 32 {
		return errShortKeySize
	}
	c.key = key

	return nil
}

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
