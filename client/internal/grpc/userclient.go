package grpc

import (
	"context"

	pb "github.com/rycln/gokeep/pkg/gen/grpc/gophkeeper"
	"github.com/rycln/gokeep/shared/models"
	"google.golang.org/grpc"
)

// AuthClient handles authentication operations via gRPC
type GophKeeperClient struct {
	client pb.GophKeeperClient // gRPC generated client interface
}

// NewUserClient creates new AuthClient instance
func NewUserClient(conn *grpc.ClientConn) *GophKeeperClient {
	return &GophKeeperClient{
		client: pb.NewGophKeeperClient(conn),
	}
}

// Register performs user registration via gRPC
func (c *GophKeeperClient) Register(ctx context.Context, req *models.UserRegReq) (*models.User, error) {
	res, err := c.client.Register(ctx, &pb.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
		Salt:     req.Salt,
	})

	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:   models.UserID(res.UserId),
		JWT:  res.Token,
		Salt: res.Salt,
	}, err
}

// Login performs user authentication via gRPC
func (c *GophKeeperClient) Login(ctx context.Context, req *models.UserLoginReq) (*models.User, error) {
	res, err := c.client.Login(ctx, &pb.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:   models.UserID(res.UserId),
		JWT:  res.Token,
		Salt: res.Salt,
	}, err
}
