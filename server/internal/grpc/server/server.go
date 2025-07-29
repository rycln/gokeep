package server

import (
	pb "github.com/rycln/gokeep/pkg/gen/grpc/gophkeeper"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	uservice userService
}

func NewGophKeeperServer(
	uservice userService,
) *UserHandler {
	return &UserHandler{
		uservice: uservice,
	}
}
