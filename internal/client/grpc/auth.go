package grpc

import (
	"context"

	pb "github.com/rycln/gokeep/internal/shared/proto/gen/gophkeeper"
)

type AuthClient struct {
	client pb.GophKeeperClient
}

func (c *AuthClient) Register(ctx context.Context, username, password string) error {
	_, err := c.client.Register(ctx, &pb.RegisterRequest{
		Username: username,
		Password: password,
	})
	return err
}
