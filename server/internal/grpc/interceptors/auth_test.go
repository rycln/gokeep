package interceptors

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rycln/gokeep/server/internal/contextkeys"
	"github.com/rycln/gokeep/server/internal/grpc/interceptors/mocks"
	"github.com/rycln/gokeep/server/internal/logger"
	"github.com/rycln/gokeep/shared/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"google.golang.org/grpc/metadata"
)

const (
	testUserID = models.UserID("550e8400-e29b-41d4-a716-446655440000")
	testToken  = "test.jwt.token"
)

func TestNewAuthInterceptor(t *testing.T) {
	t.Run("should create new interceptor", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockjwtServicer(ctrl)
		interceptor := NewAuthInterceptor(mockService)

		assert.NotNil(t, interceptor)
		assert.Equal(t, mockService, interceptor.authService)
	})
}

func TestAuthInterceptor_AuthFunc(t *testing.T) {
	observedZapCore, observedLogs := observer.New(zap.DebugLevel)
	originalLogger := logger.Log
	logger.Log = zap.New(observedZapCore)
	defer func() { logger.Log = originalLogger }()

	t.Run("successful authentication", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockjwtServicer(ctrl)
		interceptor := NewAuthInterceptor(mockService)

		md := metadata.Pairs("authorization", "bearer "+testToken)
		ctx := metadata.NewIncomingContext(context.Background(), md)

		mockService.EXPECT().
			ParseIDFromJWT(testToken).
			Return(testUserID, nil)

		newCtx, err := interceptor.AuthFunc(ctx)
		require.NoError(t, err)
		assert.Equal(t, testUserID, newCtx.Value(contextkeys.UserID))
		assert.Equal(t, 0, observedLogs.Len())
	})

	t.Run("missing auth header", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockjwtServicer(ctrl)
		interceptor := NewAuthInterceptor(mockService)

		ctx := context.Background()

		_, err := interceptor.AuthFunc(ctx)
		require.Error(t, err)

		require.Equal(t, 1, observedLogs.Len())
		log := observedLogs.All()[0]
		assert.Equal(t, "auth interceptor", log.Message)
		observedLogs.TakeAll()
	})

	t.Run("invalid token format", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockjwtServicer(ctrl)
		interceptor := NewAuthInterceptor(mockService)

		md := metadata.Pairs("authorization", testToken)
		ctx := metadata.NewIncomingContext(context.Background(), md)

		_, err := interceptor.AuthFunc(ctx)
		require.Error(t, err)

		require.Equal(t, 1, observedLogs.Len())
		log := observedLogs.All()[0]
		assert.Equal(t, "auth interceptor", log.Message)
		observedLogs.TakeAll()
	})

	t.Run("invalid JWT token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockjwtServicer(ctrl)
		interceptor := NewAuthInterceptor(mockService)

		md := metadata.Pairs("authorization", "bearer "+testToken)
		ctx := metadata.NewIncomingContext(context.Background(), md)

		testErr := errors.New("invalid token")
		mockService.EXPECT().
			ParseIDFromJWT(testToken).
			Return(models.UserID(""), testErr)

		_, err := interceptor.AuthFunc(ctx)
		require.Error(t, err)
		assert.Equal(t, testErr, err)

		require.Equal(t, 1, observedLogs.Len())
		log := observedLogs.All()[0]
		assert.Equal(t, "auth interceptor", log.Message)
		observedLogs.TakeAll()
	})
}
