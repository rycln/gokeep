package server

import (
	"context"
	"time"

	"github.com/rycln/gokeep/internal/shared/models"
	pb "github.com/rycln/gokeep/internal/shared/proto/gen/gophkeeper"
)

const timeout = time.Duration(5) * time.Second

type authServicer interface {
	CreateUser(context.Context, *models.UserAuthReq) (*models.User, error)
	UserAuth(context.Context, *models.UserAuthReq) (*models.User, error)
}

func (s *GophKeeperServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	authReq := &models.UserAuthReq{
		Username: req.Username,
		Password: req.Password,
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	user, err := s.auth.CreateUser(ctx, authReq)
	if err != nil {
		return nil, err
	}

	return &pb.AuthResponse{
		UserId:       string(user.ID),
		Token:        user.JWT,
		RefreshToken: user.RefJWT,
	}, nil
}

func (s *GophKeeperServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	authReq := &models.UserAuthReq{
		Username: req.Username,
		Password: req.Password,
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	user, err := s.auth.UserAuth(ctx, authReq)
	if err != nil {
		return nil, err
	}

	return &pb.AuthResponse{
		UserId:       string(user.ID),
		Token:        user.JWT,
		RefreshToken: user.RefJWT,
	}, nil
}
