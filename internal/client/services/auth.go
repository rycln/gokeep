package services

import (
	"context"

	"github.com/rycln/gokeep/internal/shared/models"
)

type AuthAPI interface {
	Register(context.Context, *models.UserAuthReq) (*models.User, error)
	Login(context.Context, *models.UserAuthReq) (*models.User, error)
}

type AuthService struct {
	api AuthAPI
}

func NewAuthService(api AuthAPI) *AuthService {
	return &AuthService{
		api: api,
	}
}

func (s *AuthService) UserRegister(ctx context.Context, req *models.UserAuthReq) (*models.User, error) {
	return s.api.Register(ctx, req)
}

func (s *AuthService) UserLogin(ctx context.Context, req *models.UserAuthReq) (*models.User, error) {
	return s.api.Login(ctx, req)
}
