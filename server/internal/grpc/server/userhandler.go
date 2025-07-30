package server

import (
	"context"

	pb "github.com/rycln/gokeep/pkg/gen/grpc/gophkeeper"
	"github.com/rycln/gokeep/shared/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// userService defines the required domain operations for user management
type userService interface {
	CreateUser(context.Context, *models.UserAuthReq) (*models.User, error) // User registration
	AuthUser(context.Context, *models.UserAuthReq) (*models.User, error)   // User authentication
}

// Register handles user registration requests
func (h *UserHandler) Register(
	ctx context.Context,
	req *pb.RegisterRequest,
) (*pb.AuthResponse, error) {
	authReq := &models.UserAuthReq{
		Username: req.Username,
		Password: req.Password,
	}

	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	user, err := h.uservice.CreateUser(ctx, authReq)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.AuthResponse{
		UserId: string(user.ID),
		Token:  user.JWT,
	}, nil
}

// Login handles user authentication requests
func (h *UserHandler) Login(
	ctx context.Context,
	req *pb.LoginRequest,
) (*pb.AuthResponse, error) {
	authReq := &models.UserAuthReq{
		Username: req.Username,
		Password: req.Password,
	}

	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	user, err := h.uservice.AuthUser(ctx, authReq)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.AuthResponse{
		UserId: string(user.ID),
		Token:  user.JWT,
	}, nil
}

// AuthFuncOverride provides authentication middleware hook
func (s *UserHandler) AuthFuncOverride(
	ctx context.Context,
	fullMethodName string,
) (context.Context, error) {
	return ctx, nil
}
