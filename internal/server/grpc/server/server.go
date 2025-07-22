package server

import (
	pb "github.com/rycln/gokeep/internal/shared/proto/gen/gophkeeper"
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
