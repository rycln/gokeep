package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rycln/gokeep/shared/models"
)

// Error definitions
var errNoUserID = errors.New("does not contain user id")

// JWTService handles JWT token operations
type JWTService struct {
	jwtKey string
	jwtExp time.Duration
}

// NewJWTService creates a new JWT service instance
func NewJWTService(jwtkey string, jwtExp time.Duration) *JWTService {
	return &JWTService{
		jwtKey: jwtkey,
		jwtExp: jwtExp,
	}
}

// jwtClaims contains custom JWT claims structure
type jwtClaims struct {
	jwt.RegisteredClaims               // Standard JWT claims
	UserID               models.UserID `json:"id"` // Custom user ID claim
}

// Validate implements jwt.ClaimsValidator interface
func (c jwtClaims) Validate() error {
	if c.UserID == "" {
		return errNoUserID
	}
	return nil
}

// NewJWTString generates a new signed JWT token
func (s *JWTService) NewJWTString(userID models.UserID) (string, error) {
	if userID == "" {
		return "", errNoUserID
	}

	claims := jwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.jwtExp)),
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseIDFromJWT extracts user ID from JWT token
func (s *JWTService) ParseIDFromJWT(token string) (models.UserID, error) {
	claims := &jwtClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.jwtKey), nil
	})
	if err != nil {
		return "", err
	}

	return claims.UserID, nil
}
