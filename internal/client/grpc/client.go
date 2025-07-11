package grpc

import (
	pb "github.com/rycln/gokeep/internal/shared/proto/gen/gophkeeper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	Auth *AuthClient
	conn *grpc.ClientConn
}

func New(addr string) (*GrpcClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	pbClient := pb.NewGophKeeperClient(conn)
	return &GrpcClient{
		Auth: &AuthClient{client: pbClient},
		conn: conn,
	}, nil
}
