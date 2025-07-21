package server

import (
	"context"
	"time"

	"github.com/rycln/gokeep/internal/shared/models"
	pb "github.com/rycln/gokeep/internal/shared/proto/gen/gophkeeper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const timeout = time.Duration(5) * time.Second

type authServicer interface {
	CreateUser(context.Context, *models.UserAuthReq) (*models.User, error)
	AuthUser(context.Context, *models.UserAuthReq) (*models.User, error)
}

func (h *UserHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	authReq := &models.UserAuthReq{
		Username: req.Username,
		Password: req.Password,
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	user, err := h.auth.CreateUser(ctx, authReq)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.AuthResponse{
		UserId: string(user.ID),
		Token:  user.JWT,
	}, nil
}

func (s *UserHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	authReq := &models.UserAuthReq{
		Username: req.Username,
		Password: req.Password,
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	user, err := s.auth.AuthUser(ctx, authReq)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.AuthResponse{
		UserId: string(user.ID),
		Token:  user.JWT,
	}, nil
}

func (s *UserHandler) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	return ctx, nil
}
