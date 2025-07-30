// Package server implements the gRPC service handlers.
package server

import (
	"time"

	pb "github.com/rycln/gokeep/pkg/gen/grpc/gophkeeper"
)

// UserHandler implements the gRPC UserService server interface.
// Wraps user domain operations with gRPC protocol handling.
type UserHandler struct {
	pb.UnimplementedUserServiceServer
	uservice userService
	timeout  time.Duration
}

// NewGophKeeperServer constructs a new gRPC server instance with required dependencies.
// Returns configured server ready for registration with gRPC
func NewGophKeeperServer(uservice userService, timeout time.Duration) *UserHandler {
	return &UserHandler{
		uservice: uservice,
		timeout:  timeout,
	}
}
