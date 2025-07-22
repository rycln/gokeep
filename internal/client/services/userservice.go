package services

import (
	"context"

	"github.com/rycln/gokeep/internal/shared/models"
)

type AuthAPI interface {
	Register(context.Context, *models.UserAuthReq) (*models.User, error)
	Login(context.Context, *models.UserAuthReq) (*models.User, error)
}

type UserService struct {
	api AuthAPI
}

func NewAuthService(api AuthAPI) *UserService {
	return &UserService{
		api: api,
	}
}

func (s *UserService) UserRegister(ctx context.Context, req *models.UserAuthReq) (*models.User, error) {
	return s.api.Register(ctx, req)
}

func (s *UserService) UserLogin(ctx context.Context, req *models.UserAuthReq) (*models.User, error) {
	return s.api.Login(ctx, req)
}
