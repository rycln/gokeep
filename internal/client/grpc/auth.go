package grpc

import (
	"context"

	"github.com/rycln/gokeep/internal/shared/models"
	pb "github.com/rycln/gokeep/internal/shared/proto/gen/gophkeeper"
	"google.golang.org/grpc"
)

type AuthClient struct {
	client pb.UserServiceClient
}

func NewAuthClient(conn *grpc.ClientConn) *AuthClient {
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
		ID:     models.UserID(res.UserId),
		JWT:    res.Token,
		RefJWT: res.RefreshToken,
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
		ID:     models.UserID(res.UserId),
		JWT:    res.Token,
		RefJWT: res.RefreshToken,
	}, err
}
