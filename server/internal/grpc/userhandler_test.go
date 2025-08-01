package grpc

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

func TestGophKeeperServer_Register(t *testing.T) {
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

		mockUser := mocks.NewMockuserService(ctrl)
		mockSync := mocks.NewMocksyncService(ctrl)
		mockAuth := mocks.NewMockauthProvider(ctrl)
		handler := NewGophKeeperServer(mockUser, mockSync, mockAuth, testTimeout)

		expectedUser := &models.User{
			ID:   models.UserID(testUserID),
			JWT:  testJWT,
			Salt: testSalt,
		}

		mockUser.EXPECT().
			CreateUser(gomock.Any(), expectedAuthReq).
			DoAndReturn(func(ctx context.Context, _ *models.UserRegReq) (*models.User, error) {
				_, ok := ctx.Deadline()
				assert.True(t, ok, "context should have deadline")
				return expectedUser, nil
			})

		resp, err := handler.Register(context.Background(), testReq)
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, testUserID, resp.UserId)
		assert.Equal(t, testJWT, resp.Token)
		assert.Equal(t, testSalt, resp.Salt)
	})

	t.Run("service error returns grpc error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUser := mocks.NewMockuserService(ctrl)
		mockSync := mocks.NewMocksyncService(ctrl)
		mockAuth := mocks.NewMockauthProvider(ctrl)
		handler := NewGophKeeperServer(mockUser, mockSync, mockAuth, testTimeout)

		testErr := errors.New("test error")
		mockUser.EXPECT().
			CreateUser(gomock.Any(), expectedAuthReq).
			Return(nil, testErr)

		resp, err := handler.Register(context.Background(), testReq)
		require.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
		assert.Contains(t, err.Error(), testErr.Error())
	})
}

func TestGophKeeperServer_Login(t *testing.T) {
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

		mockUser := mocks.NewMockuserService(ctrl)
		mockSync := mocks.NewMocksyncService(ctrl)
		mockAuth := mocks.NewMockauthProvider(ctrl)
		handler := NewGophKeeperServer(mockUser, mockSync, mockAuth, testTimeout)

		expectedUser := &models.User{
			ID:   models.UserID(testUserID),
			JWT:  testJWT,
			Salt: testSalt,
		}

		mockUser.EXPECT().
			AuthUser(gomock.Any(), expectedAuthReq).
			DoAndReturn(func(ctx context.Context, _ *models.UserLoginReq) (*models.User, error) {
				_, ok := ctx.Deadline()
				assert.True(t, ok, "context should have deadline")
				return expectedUser, nil
			})

		resp, err := handler.Login(context.Background(), testReq)
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, testUserID, resp.UserId)
		assert.Equal(t, testJWT, resp.Token)
		assert.Equal(t, testSalt, resp.Salt)
	})

	t.Run("service error returns grpc error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUser := mocks.NewMockuserService(ctrl)
		mockSync := mocks.NewMocksyncService(ctrl)
		mockAuth := mocks.NewMockauthProvider(ctrl)
		handler := NewGophKeeperServer(mockUser, mockSync, mockAuth, testTimeout)

		testErr := errors.New("test error")
		mockUser.EXPECT().
			AuthUser(gomock.Any(), expectedAuthReq).
			Return(nil, testErr)

		resp, err := handler.Login(context.Background(), testReq)
		require.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
		assert.Contains(t, err.Error(), testErr.Error())
	})
}

func TestGophKeeperServer_AuthFuncOverride(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUser := mocks.NewMockuserService(ctrl)
	mockSync := mocks.NewMocksyncService(ctrl)
	mockAuth := mocks.NewMockauthProvider(ctrl)
	server := NewGophKeeperServer(mockUser, mockSync, mockAuth, testTimeout)

	t.Run("should bypass auth for Register method", func(t *testing.T) {
		ctx := context.Background()
		fullMethodName := "/gophkeeper.GophKeeper/Register"

		resultCtx, err := server.AuthFuncOverride(ctx, fullMethodName)

		assert.NoError(t, err)
		assert.Equal(t, ctx, resultCtx)
	})

	t.Run("should bypass auth for Login method", func(t *testing.T) {
		ctx := context.Background()
		fullMethodName := "/gophkeeper.GophKeeper/Login"

		resultCtx, err := server.AuthFuncOverride(ctx, fullMethodName)

		assert.NoError(t, err)
		assert.Equal(t, ctx, resultCtx)
	})

	t.Run("should handle case-insensitive method names", func(t *testing.T) {
		testCases := []struct {
			name       string
			method     string
			shouldAuth bool
		}{
			{"other method", "/service/Other", true},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				ctx := context.Background()

				if tc.shouldAuth {
					mockAuth.EXPECT().
						AuthFunc(ctx).
						Return(ctx, nil)
				}

				_, err := server.AuthFuncOverride(ctx, tc.method)

				assert.NoError(t, err)
			})
		}
	})
}
