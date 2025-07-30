package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/rycln/gokeep/pkg/gen/grpc/gophkeeper"
	"github.com/rycln/gokeep/shared/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

const (
	testUserID = "550e8400-e29b-41d4-a716-446655440000"
	testToken  = "test.jwt.token"
	testUser   = "testuser"
	testPass   = "testpass"
)

type mockUserServiceClient struct {
	gophkeeper.UserServiceClient
	registerFunc func(ctx context.Context, in *gophkeeper.RegisterRequest, opts ...grpc.CallOption) (*gophkeeper.AuthResponse, error)
	loginFunc    func(ctx context.Context, in *gophkeeper.LoginRequest, opts ...grpc.CallOption) (*gophkeeper.AuthResponse, error)
}

func (m *mockUserServiceClient) Register(ctx context.Context, in *gophkeeper.RegisterRequest, opts ...grpc.CallOption) (*gophkeeper.AuthResponse, error) {
	return m.registerFunc(ctx, in, opts...)
}

func (m *mockUserServiceClient) Login(ctx context.Context, in *gophkeeper.LoginRequest, opts ...grpc.CallOption) (*gophkeeper.AuthResponse, error) {
	return m.loginFunc(ctx, in, opts...)
}

func TestNewUserClient(t *testing.T) {
	t.Run("should create new auth client", func(t *testing.T) {
		conn := &grpc.ClientConn{}
		client := NewUserClient(conn)

		assert.NotNil(t, client)
		assert.NotNil(t, client.client)
	})
}

func TestAuthClient_Register(t *testing.T) {
	ctx := context.Background()
	testReq := &models.UserAuthReq{
		Username: testUser,
		Password: testPass,
	}

	t.Run("successful registration", func(t *testing.T) {
		mockClient := &mockUserServiceClient{
			registerFunc: func(ctx context.Context, in *gophkeeper.RegisterRequest, opts ...grpc.CallOption) (*gophkeeper.AuthResponse, error) {
				assert.Equal(t, testUser, in.Username)
				assert.Equal(t, testPass, in.Password)
				return &gophkeeper.AuthResponse{
					UserId: testUserID,
					Token:  testToken,
				}, nil
			},
		}

		authClient := &AuthClient{client: mockClient}
		user, err := authClient.Register(ctx, testReq)

		require.NoError(t, err)
		assert.Equal(t, models.UserID(testUserID), user.ID)
		assert.Equal(t, testToken, user.JWT)
	})

	t.Run("registration error", func(t *testing.T) {
		expectedErr := errors.New("registration failed")
		mockClient := &mockUserServiceClient{
			registerFunc: func(ctx context.Context, in *gophkeeper.RegisterRequest, opts ...grpc.CallOption) (*gophkeeper.AuthResponse, error) {
				return nil, expectedErr
			},
		}

		authClient := &AuthClient{client: mockClient}
		_, err := authClient.Register(ctx, testReq)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestAuthClient_Login(t *testing.T) {
	ctx := context.Background()
	testReq := &models.UserAuthReq{
		Username: testUser,
		Password: testPass,
	}

	t.Run("successful login", func(t *testing.T) {
		mockClient := &mockUserServiceClient{
			loginFunc: func(ctx context.Context, in *gophkeeper.LoginRequest, opts ...grpc.CallOption) (*gophkeeper.AuthResponse, error) {
				assert.Equal(t, testUser, in.Username)
				assert.Equal(t, testPass, in.Password)
				return &gophkeeper.AuthResponse{
					UserId: testUserID,
					Token:  testToken,
				}, nil
			},
		}

		authClient := &AuthClient{client: mockClient}
		user, err := authClient.Login(ctx, testReq)

		require.NoError(t, err)
		assert.Equal(t, models.UserID(testUserID), user.ID)
		assert.Equal(t, testToken, user.JWT)
	})

	t.Run("login error", func(t *testing.T) {
		expectedErr := errors.New("login failed")
		mockClient := &mockUserServiceClient{
			loginFunc: func(ctx context.Context, in *gophkeeper.LoginRequest, opts ...grpc.CallOption) (*gophkeeper.AuthResponse, error) {
				return nil, expectedErr
			},
		}

		authClient := &AuthClient{client: mockClient}
		_, err := authClient.Login(ctx, testReq)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}
