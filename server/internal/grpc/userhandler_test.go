package server

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/rycln/gokeep/pkg/gen/grpc/gophkeeper"
	"github.com/rycln/gokeep/server/internal/grpc/mocks"
	"github.com/rycln/gokeep/shared/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	testUserID  = "550e8400-e29b-41d4-a716-446655440000"
	testJWT     = "test.jwt.token"
	testSalt    = "salt"
	testTimeout = 5 * time.Second
)

func TestUserHandler_Register(t *testing.T) {
	testReq := &gophkeeper.RegisterRequest{
		Username: "testuser",
		Password: "testpass",
		Salt:     testSalt,
	}

	expectedAuthReq := &models.UserRegReq{
		Username: testReq.Username,
		Password: testReq.Password,
		Salt:     testReq.Salt,
	}

	t.Run("successful registration", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockuserService(ctrl)
		handler := NewGophKeeperServer(mockService, testTimeout)

		expectedUser := &models.User{
			ID:  models.UserID(testUserID),
			JWT: testJWT,
		}

		mockService.EXPECT().
			CreateUser(gomock.Any(), expectedAuthReq).
			Return(expectedUser, nil)

		resp, err := handler.Register(context.Background(), testReq)
		require.NoError(t, err)
		assert.Equal(t, testUserID, resp.UserId)
		assert.Equal(t, testJWT, resp.Token)
	})

	t.Run("service error returns grpc error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockuserService(ctrl)
		handler := NewGophKeeperServer(mockService, testTimeout)

		testErr := errors.New("test error")
		mockService.EXPECT().
			CreateUser(gomock.Any(), expectedAuthReq).
			Return(nil, testErr)

		_, err := handler.Register(context.Background(), testReq)
		require.Error(t, err)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
		assert.Contains(t, err.Error(), testErr.Error())
	})
}

func TestUserHandler_Login(t *testing.T) {
	testReq := &gophkeeper.LoginRequest{
		Username: "testuser",
		Password: "testpass",
	}

	expectedAuthReq := &models.UserLoginReq{
		Username: testReq.Username,
		Password: testReq.Password,
	}

	t.Run("successful login", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockuserService(ctrl)
		handler := NewGophKeeperServer(mockService, testTimeout)

		expectedUser := &models.User{
			ID:   models.UserID(testUserID),
			JWT:  testJWT,
			Salt: testSalt,
		}

		mockService.EXPECT().
			AuthUser(gomock.Any(), expectedAuthReq).
			Return(expectedUser, nil)

		resp, err := handler.Login(context.Background(), testReq)
		require.NoError(t, err)
		assert.Equal(t, testUserID, resp.UserId)
		assert.Equal(t, testJWT, resp.Token)
	})

	t.Run("service error returns grpc error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockuserService(ctrl)
		handler := NewGophKeeperServer(mockService, testTimeout)

		testErr := errors.New("test error")
		mockService.EXPECT().
			AuthUser(gomock.Any(), expectedAuthReq).
			Return(nil, testErr)

		_, err := handler.Login(context.Background(), testReq)
		require.Error(t, err)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
		assert.Contains(t, err.Error(), testErr.Error())
	})
}

func TestUserHandler_AuthFuncOverride(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockuserService(ctrl)
	handler := NewGophKeeperServer(mockService, testTimeout)

	t.Run("always returns nil", func(t *testing.T) {
		ctx := context.Background()
		newCtx, err := handler.AuthFuncOverride(ctx, "test/method")
		assert.NoError(t, err)
		assert.Equal(t, ctx, newCtx)
	})
}
