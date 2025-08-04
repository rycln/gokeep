package crypto

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAESCrypter(t *testing.T) {
	t.Run("should create new crypter with empty key", func(t *testing.T) {
		c := NewAESCrypter()
		assert.NotNil(t, c)
		assert.Nil(t, c.key)
	})
}

func TestAESCrypter_SetKey(t *testing.T) {
	t.Run("should accept valid 32-byte key", func(t *testing.T) {
		c := NewAESCrypter()
		key := make([]byte, 32)
		_, err := rand.Read(key)
		require.NoError(t, err)

		err = c.SetKey(key)
		assert.NoError(t, err)
		assert.Equal(t, key, c.key)
	})

	t.Run("should reject short key", func(t *testing.T) {
		c := NewAESCrypter()
		err := c.SetKey(make([]byte, 16))
		assert.Error(t, err)
		assert.Equal(t, errShortKeySize, err)
		assert.Nil(t, c.key)
	})

	t.Run("should reject long key", func(t *testing.T) {
		c := NewAESCrypter()
		err := c.SetKey(make([]byte, 64))
		assert.Error(t, err)
		assert.Equal(t, errShortKeySize, err)
		assert.Nil(t, c.key)
	})
}

func TestAESCrypter_Encrypt(t *testing.T) {
	t.Run("should encrypt data with valid key", func(t *testing.T) {
		c := NewAESCrypter()
		key := make([]byte, 32)
		_, err := rand.Read(key)
		require.NoError(t, err)
		require.NoError(t, c.SetKey(key))

		data := []byte("secret data")
		encrypted, err := c.Encrypt(data)
		require.NoError(t, err)
		assert.NotEmpty(t, encrypted)
		assert.NotEqual(t, data, encrypted)
	})

	t.Run("should fail without key", func(t *testing.T) {
		c := NewAESCrypter()
		_, err := c.Encrypt([]byte("data"))
		assert.Error(t, err)
		assert.Equal(t, errNoKey, err)
	})

	t.Run("should encrypt empty data", func(t *testing.T) {
		c := NewAESCrypter()
		key := make([]byte, 32)
		_, err := rand.Read(key)
		require.NoError(t, err)
		require.NoError(t, c.SetKey(key))

		encrypted, err := c.Encrypt([]byte{})
		require.NoError(t, err)
		assert.NotEmpty(t, encrypted)
	})
}

func TestAESCrypter_Decrypt(t *testing.T) {
	t.Run("should decrypt valid ciphertext", func(t *testing.T) {
		c := NewAESCrypter()
		key := make([]byte, 32)
		_, err := rand.Read(key)
		require.NoError(t, err)
		require.NoError(t, c.SetKey(key))

		original := []byte("secret message")
		encrypted, err := c.Encrypt(original)
		require.NoError(t, err)

		decrypted, err := c.Decrypt(encrypted)
		require.NoError(t, err)
		assert.Equal(t, original, decrypted)
	})

	t.Run("should fail without key", func(t *testing.T) {
		c := NewAESCrypter()
		_, err := c.Decrypt([]byte("invalid"))
		assert.Error(t, err)
		assert.Equal(t, errNoKey, err)
	})

	t.Run("should fail with short ciphertext", func(t *testing.T) {
		c := NewAESCrypter()
		key := make([]byte, 32)
		_, err := rand.Read(key)
		require.NoError(t, err)
		require.NoError(t, c.SetKey(key))

		_, err = c.Decrypt([]byte("too short"))
		assert.Error(t, err)
		assert.Equal(t, errShortCrypted, err)
	})

	t.Run("should fail with tampered ciphertext", func(t *testing.T) {
		c := NewAESCrypter()
		key := make([]byte, 32)
		_, err := rand.Read(key)
		require.NoError(t, err)
		require.NoError(t, c.SetKey(key))

		original := []byte("secret message")
		encrypted, err := c.Encrypt(original)
		require.NoError(t, err)

		// Tamper with the ciphertext
		encrypted[10] ^= 0xFF

		_, err = c.Decrypt(encrypted)
		assert.Error(t, err)
	})
}

func TestAESCrypter_EncryptDecrypt(t *testing.T) {
	t.Run("should roundtrip various data sizes", func(t *testing.T) {
		c := NewAESCrypter()
		key := make([]byte, 32)
		_, err := rand.Read(key)
		require.NoError(t, err)
		require.NoError(t, c.SetKey(key))

		testCases := []struct {
			name string
			data []byte
		}{
			{"small", []byte("a")},
			{"medium", []byte("some test data")},
			{"large", make([]byte, 1024)},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				encrypted, err := c.Encrypt(tc.data)
				require.NoError(t, err)

				decrypted, err := c.Decrypt(encrypted)
				require.NoError(t, err)

				assert.Equal(t, tc.data, decrypted)
			})
		}
	})

	t.Run("should fail with wrong key", func(t *testing.T) {
		// Create first crypter and encrypt
		c1 := NewAESCrypter()
		key1 := make([]byte, 32)
		_, err := rand.Read(key1)
		require.NoError(t, err)
		require.NoError(t, c1.SetKey(key1))

		data := []byte("secret data")
		encrypted, err := c1.Encrypt(data)
		require.NoError(t, err)

		// Create second crypter with different key
		c2 := NewAESCrypter()
		key2 := make([]byte, 32)
		_, err = rand.Read(key2)
		require.NoError(t, err)
		require.NoError(t, c2.SetKey(key2))

		_, err = c2.Decrypt(encrypted)
		assert.Error(t, err)
	})
}
