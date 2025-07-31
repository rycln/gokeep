package services

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewKeyService(t *testing.T) {
	t.Run("should create new KeyService instance", func(t *testing.T) {
		service := NewKeyService()
		assert.NotNil(t, service)
	})
}

func TestGenerateSalt(t *testing.T) {
	t.Run("should generate valid salt", func(t *testing.T) {
		service := NewKeyService()
		salt, err := service.GenerateSalt()

		assert.NoError(t, err)
		assert.Equal(t, saltLength, len(salt))

		anotherSalt, _ := service.GenerateSalt()
		assert.NotEqual(t, salt, anotherSalt)
	})
}

func TestDeriveKeyFromPasswordAndSalt(t *testing.T) {
	t.Run("should derive consistent key from same inputs", func(t *testing.T) {
		service := NewKeyService()
		password := "securePassword123"
		salt := []byte("testSalt123456789") // 16 bytes

		key1 := service.DeriveKeyFromPasswordAndSalt(password, salt)
		key2 := service.DeriveKeyFromPasswordAndSalt(password, salt)

		assert.Equal(t, keyLength, len(key1))
		assert.Equal(t, key1, key2)
	})

	t.Run("should produce different keys for different passwords", func(t *testing.T) {
		service := NewKeyService()
		salt := []byte("testSalt123456789")

		key1 := service.DeriveKeyFromPasswordAndSalt("password1", salt)
		key2 := service.DeriveKeyFromPasswordAndSalt("password2", salt)

		assert.NotEqual(t, key1, key2)
	})

	t.Run("should produce different keys for different salts", func(t *testing.T) {
		service := NewKeyService()
		password := "password"

		key1 := service.DeriveKeyFromPasswordAndSalt(password, []byte("salt1_123456789"))
		key2 := service.DeriveKeyFromPasswordAndSalt(password, []byte("salt2_123456789"))

		assert.NotEqual(t, key1, key2)
	})
}

func TestDecodeSalt(t *testing.T) {
	t.Run("should decode valid base64 salt", func(t *testing.T) {
		service := NewKeyService()
		original := []byte("testSalt12345678") // 16 bytes
		encoded := base64.StdEncoding.EncodeToString(original)

		decoded, err := service.DecodeSalt(encoded)

		assert.NoError(t, err)
		assert.Equal(t, original, decoded)
	})

	t.Run("should return error for invalid base64", func(t *testing.T) {
		service := NewKeyService()
		_, err := service.DecodeSalt("not valid base64!")
		assert.Error(t, err)
	})

	t.Run("should return error for incorrect salt length", func(t *testing.T) {
		service := NewKeyService()
		shortSalt := []byte("short")
		encoded := base64.StdEncoding.EncodeToString(shortSalt)

		_, err := service.DecodeSalt(encoded)
		assert.Equal(t, errInvalidSaltLength, err)
	})
}

func TestEncodeSalt(t *testing.T) {
	t.Run("should encode salt to base64", func(t *testing.T) {
		service := NewKeyService()
		salt := []byte("testSalt123456789") // 16 bytes

		encoded := service.EncodeSalt(salt)
		decoded, _ := base64.StdEncoding.DecodeString(encoded)

		assert.Equal(t, salt, decoded)
	})

	t.Run("should handle empty salt", func(t *testing.T) {
		service := NewKeyService()
		encoded := service.EncodeSalt([]byte{})
		assert.Empty(t, encoded)
	})
}
