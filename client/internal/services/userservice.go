package services

import (
	"context"

	"github.com/rycln/gokeep/shared/models"
)

// AuthAPI defines the interface for authentication operations
type AuthAPI interface {
	Register(context.Context, *models.UserAuthReq) (*models.User, error)
	Login(context.Context, *models.UserAuthReq) (*models.User, error)
}

// UserService handles user authentication business logic
type UserService struct {
	api AuthAPI // Authentication API implementation
}

// NewAuthService creates a new UserService instance
func NewAuthService(api AuthAPI) *UserService {
	return &UserService{
		api: api,
	}
}

// UserRegister handles user registration flow
func (s *UserService) UserRegister(ctx context.Context, req *models.UserAuthReq) (*models.User, error) {
	return s.api.Register(ctx, req)
}

// UserLogin handles user authentication flow
func (s *UserService) UserLogin(ctx context.Context, req *models.UserAuthReq) (*models.User, error) {
	return s.api.Login(ctx, req)
}
