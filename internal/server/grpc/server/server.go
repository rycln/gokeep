package server

import (
	pb "github.com/rycln/gokeep/internal/shared/proto/gen/gophkeeper"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	auth authServicer
}

func NewGophKeeperServer(
	auth authServicer,
) *UserHandler {
	return &UserHandler{
		auth: auth,
	}
}
