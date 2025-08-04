// Package grpc implements the gRPC service.
package grpc

import (
	"time"

	pb "github.com/rycln/gokeep/pkg/gen/grpc/gophkeeper"
)

// GophKeeperServer implements the gRPC server interface.
type GophKeeperServer struct {
	pb.UnimplementedGophKeeperServer
	user    userService
	sync    syncService
	auth    authProvider
	timeout time.Duration
}

// NewGophKeeperServer constructs a new gRPC server instance with required dependencies.
// Returns configured server ready for registration with gRPC
func NewGophKeeperServer(
	user userService,
	sync syncService,
	auth authProvider,
	timeout time.Duration,
) *GophKeeperServer {
	return &GophKeeperServer{
		user:    user,
		sync:    sync,
		auth:    auth,
		timeout: timeout,
	}
}
