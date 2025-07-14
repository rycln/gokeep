package server

import (
	pb "github.com/rycln/gokeep/internal/shared/proto/gen/gophkeeper"
)

type GophKeeperServer struct {
	pb.UnimplementedGophKeeperServer

	auth authServicer
}

func NewGophKeeperServer(
	auth authServicer,
) *GophKeeperServer {
	return &GophKeeperServer{
		auth: auth,
	}
}
