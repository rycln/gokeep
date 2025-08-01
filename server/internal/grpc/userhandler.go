package grpc

import (
	"context"
	"strings"

	pb "github.com/rycln/gokeep/pkg/gen/grpc/gophkeeper"
	"github.com/rycln/gokeep/shared/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

// userService defines the required domain operations for user management
type userService interface {
	CreateUser(context.Context, *models.UserRegReq) (*models.User, error) // User registration
	AuthUser(context.Context, *models.UserLoginReq) (*models.User, error) // User authentication
}

// authProvider defines authentication middleware function
type authProvider interface {
	AuthFunc(context.Context) (context.Context, error)
}

// Register handles user registration requests
func (h *GophKeeperServer) Register(
	ctx context.Context,
	req *pb.RegisterRequest,
) (*pb.AuthResponse, error) {
	authReq := &models.UserRegReq{
		Username: req.Username,
		Password: req.Password,
		Salt:     req.Salt,
	}

	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	user, err := h.user.CreateUser(ctx, authReq)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.AuthResponse{
		UserId: string(user.ID),
		Token:  user.JWT,
		Salt:   req.Salt,
	}, nil
}

// Login handles user authentication requests
func (h *GophKeeperServer) Login(
	ctx context.Context,
	req *pb.LoginRequest,
) (*pb.AuthResponse, error) {
	authReq := &models.UserLoginReq{
		Username: req.Username,
		Password: req.Password,
	}

	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	user, err := h.user.AuthUser(ctx, authReq)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.AuthResponse{
		UserId: string(user.ID),
		Token:  user.JWT,
		Salt:   user.Salt,
	}, nil
}

// AuthFuncOverride provides authentication middleware hook
func (s *GophKeeperServer) AuthFuncOverride(
	ctx context.Context,
	fullMethodName string,
) (context.Context, error) {
	if strings.Contains(fullMethodName, "Register") || strings.Contains(fullMethodName, "Login") {
		return ctx, nil
	}

	return s.auth.AuthFunc(ctx)
}
