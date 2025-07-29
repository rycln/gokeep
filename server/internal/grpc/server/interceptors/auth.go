package interceptors

import (
	"context"

	"github.com/rycln/gokeep/server/internal/contextkeys"
	"github.com/rycln/gokeep/server/internal/logger"
	"github.com/rycln/gokeep/shared/models"
	"go.uber.org/zap"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
)

// authServicer defines the interface for authentication operations.
// Implementations should handle both JWT generation and parsing.
type jwtServicer interface {
	// ParseIDFromAuthHeader extracts user ID from a JWT authorization header.
	ParseIDFromJWT(string) (models.UserID, error)
}

// AuthInterceptor implements gRPC unary server interceptor for authentication.
// It handles both existing JWT validation and new user registration.
type AuthInterceptor struct {
	authService jwtServicer // Service handling JWT operations
}

// NewAuthInterceptor creates a new AuthInterceptor instance.
func NewAuthInterceptor(authService jwtServicer) *AuthInterceptor {
	return &AuthInterceptor{
		authService: authService,
	}
}

// AuthFunc performs authentication/authorization for gRPC requests.
func (i *AuthInterceptor) AuthFunc(ctx context.Context) (context.Context, error) {
	token, err := auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		logger.Log.Debug("auth interceptor", zap.Error(err))
		return nil, err
	}

	uid, err := i.authService.ParseIDFromJWT(token)
	if err != nil {
		logger.Log.Debug("auth interceptor", zap.Error(err))
		return nil, err
	}

	return context.WithValue(ctx, contextkeys.UserID, uid), nil
}
