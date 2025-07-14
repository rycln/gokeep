package services

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rycln/gokeep/internal/shared/models"
)

var errNoUserID = errors.New("does not contain user id")

type JWTService struct {
	jwtKey string
	jwtExp time.Duration
}

func NewJWTService(jwtkey string, jwtExp time.Duration) *JWTService {
	return &JWTService{
		jwtKey: jwtkey,
		jwtExp: jwtExp,
	}
}

type jwtClaims struct {
	jwt.RegisteredClaims
	UserID models.UserID `json:"id"`
}

func (c jwtClaims) Validate() error {
	if c.UserID == "" {
		return errNoUserID
	}
	return nil
}

func (s *JWTService) NewJWTString(userID models.UserID) (string, error) {
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

func (s *JWTService) ParseIDFromJWTHeader(header string) (models.UserID, error) {
	tokenString := strings.TrimPrefix(header, "Bearer")
	tokenString = strings.TrimSpace(tokenString)

	claims := &jwtClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.jwtKey), nil
	})
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}
