package grpc

import (
	"context"

	pb "github.com/rycln/gokeep/pkg/gen/grpc/gophkeeper"
	"github.com/rycln/gokeep/shared/models"
	"google.golang.org/grpc"
)

type AuthClient struct {
	client pb.UserServiceClient
}

func NewUserClient(conn *grpc.ClientConn) *AuthClient {
	return &AuthClient{
		client: pb.NewUserServiceClient(conn),
	}
}

func (c *AuthClient) Register(ctx context.Context, req *models.UserAuthReq) (*models.User, error) {
	res, err := c.client.Register(ctx, &pb.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:  models.UserID(res.UserId),
		JWT: res.Token,
	}, err
}

func (c *AuthClient) Login(ctx context.Context, req *models.UserAuthReq) (*models.User, error) {
	res, err := c.client.Login(ctx, &pb.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:  models.UserID(res.UserId),
		JWT: res.Token,
	}, err
}
