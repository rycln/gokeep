package services

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testKey = "secret_key"
	testExp = 5 * time.Second
)

func TestNewJWTService(t *testing.T) {
	t.Run("should create new JWT service", func(t *testing.T) {
		service := NewJWTService(testKey, testExp)
		assert.NotNil(t, service)
		assert.Equal(t, testKey, service.jwtKey)
		assert.Equal(t, testExp, service.jwtExp)
	})
}

func TestJWTService_NewJWTString(t *testing.T) {
	service := NewJWTService(testKey, testExp)

	t.Run("successful token generation", func(t *testing.T) {
		token, err := service.NewJWTString(testUserID)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// Verify the token can be parsed back
		userID, err := service.ParseIDFromJWT(token)
		assert.NoError(t, err)
		assert.Equal(t, testUserID, userID)
	})

	t.Run("empty user ID should return error", func(t *testing.T) {
		token, err := service.NewJWTString("")
		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Equal(t, "does not contain user id", err.Error())
	})
}

func TestJWTService_ParseIDFromJWT(t *testing.T) {
	service := NewJWTService(testKey, testExp)

	t.Run("successful token parsing", func(t *testing.T) {
		token, err := service.NewJWTString(testUserID)
		require.NoError(t, err)

		userID, err := service.ParseIDFromJWT(token)
		assert.NoError(t, err)
		assert.Equal(t, testUserID, userID)
	})

	t.Run("invalid token should fail", func(t *testing.T) {
		_, err := service.ParseIDFromJWT("invalid.token.string")
		assert.Error(t, err)
	})

	t.Run("expired token should fail", func(t *testing.T) {
		// Create service with negative expiration to generate expired token
		expiredService := NewJWTService(testKey, -time.Second)
		token, err := expiredService.NewJWTString(testUserID)
		require.NoError(t, err)

		_, err = service.ParseIDFromJWT(token)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token is expired")
	})

	t.Run("token without user ID should fail", func(t *testing.T) {
		claims := jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(testExp)),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(testKey))
		require.NoError(t, err)

		_, err = service.ParseIDFromJWT(tokenString)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not contain user id")
	})

	t.Run("wrong signing key should fail", func(t *testing.T) {
		wrongService := NewJWTService("wrong_key", testExp)
		token, err := wrongService.NewJWTString(testUserID)
		require.NoError(t, err)

		_, err = service.ParseIDFromJWT(token)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "signature is invalid")
	})
}

func TestJWTClaims_Validate(t *testing.T) {
	t.Run("valid claims", func(t *testing.T) {
		claims := jwtClaims{
			UserID: testUserID,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(testExp)),
			},
		}
		assert.NoError(t, claims.Validate())
	})

	t.Run("empty user ID should fail", func(t *testing.T) {
		claims := jwtClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(testExp)),
			},
		}
		err := claims.Validate()
		assert.Error(t, err)
		assert.Equal(t, "does not contain user id", err.Error())
	})
}
