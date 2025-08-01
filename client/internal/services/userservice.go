package services

import (
	"context"

	"github.com/rycln/gokeep/shared/models"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

// authAPI defines the interface for authentication operations
type authAPI interface {
	Register(context.Context, *models.UserRegReq) (*models.User, error)
	Login(context.Context, *models.UserLoginReq) (*models.User, error)
}

// UserService handles user authentication business logic
type UserService struct {
	api authAPI // Authentication API implementation
}

// NewAuthService creates a new UserService instance
func NewAuthService(api authAPI) *UserService {
	return &UserService{
		api: api,
	}
}

// UserRegister handles user registration flow
func (s *UserService) UserRegister(ctx context.Context, req *models.UserRegReq) (*models.User, error) {
	return s.api.Register(ctx, req)
}

// UserLogin handles user authentication flow
func (s *UserService) UserLogin(ctx context.Context, req *models.UserLoginReq) (*models.User, error) {
	return s.api.Login(ctx, req)
}
