package grpc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rycln/gokeep/pkg/gen/grpc/gophkeeper"
	"github.com/rycln/gokeep/shared/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	testUserID   = "550e8400-e29b-41d4-a716-446655440000"
	testToken    = "test.jwt.token"
	testUser     = "testuser"
	testPass     = "testpass"
	testSalt     = "salt"
	testItemID   = "550e8400-e29b-41d4-a716-446655440001"
	testMetadata = "{}"
)

type mockGophKeeperClient struct {
	gophkeeper.GophKeeperClient
	registerFunc func(ctx context.Context, in *gophkeeper.RegisterRequest, opts ...grpc.CallOption) (*gophkeeper.AuthResponse, error)
	loginFunc    func(ctx context.Context, in *gophkeeper.LoginRequest, opts ...grpc.CallOption) (*gophkeeper.AuthResponse, error)
	syncFunc     func(ctx context.Context, in *gophkeeper.SyncRequest, opts ...grpc.CallOption) (*gophkeeper.SyncResponse, error)
}

func (m *mockGophKeeperClient) Register(ctx context.Context, in *gophkeeper.RegisterRequest, opts ...grpc.CallOption) (*gophkeeper.AuthResponse, error) {
	return m.registerFunc(ctx, in, opts...)
}

func (m *mockGophKeeperClient) Login(ctx context.Context, in *gophkeeper.LoginRequest, opts ...grpc.CallOption) (*gophkeeper.AuthResponse, error) {
	return m.loginFunc(ctx, in, opts...)
}

func (m *mockGophKeeperClient) Sync(ctx context.Context, in *gophkeeper.SyncRequest, opts ...grpc.CallOption) (*gophkeeper.SyncResponse, error) {
	return m.syncFunc(ctx, in, opts...)
}

func TestNewGophKeeperClient(t *testing.T) {
	t.Run("should create new client", func(t *testing.T) {
		conn := &grpc.ClientConn{}
		client := NewGophKeeperClient(conn)

		assert.NotNil(t, client)
		assert.NotNil(t, client.client)
	})
}

func TestGophKeeperClient_Register(t *testing.T) {
	ctx := context.Background()
	testReq := &models.UserRegReq{
		Username: testUser,
		Password: testPass,
		Salt:     testSalt,
	}

	t.Run("successful registration", func(t *testing.T) {
		mockClient := &mockGophKeeperClient{
			registerFunc: func(ctx context.Context, in *gophkeeper.RegisterRequest, opts ...grpc.CallOption) (*gophkeeper.AuthResponse, error) {
				assert.Equal(t, testUser, in.Username)
				assert.Equal(t, testPass, in.Password)
				assert.Equal(t, testSalt, in.Salt)
				return &gophkeeper.AuthResponse{
					UserId: testUserID,
					Token:  testToken,
					Salt:   testSalt,
				}, nil
			},
		}

		client := &GophKeeperClient{client: mockClient}
		user, err := client.Register(ctx, testReq)

		require.NoError(t, err)
		assert.Equal(t, models.UserID(testUserID), user.ID)
		assert.Equal(t, testToken, user.JWT)
		assert.Equal(t, testSalt, user.Salt)
	})

	t.Run("registration error", func(t *testing.T) {
		expectedErr := errors.New("registration failed")
		mockClient := &mockGophKeeperClient{
			registerFunc: func(ctx context.Context, in *gophkeeper.RegisterRequest, opts ...grpc.CallOption) (*gophkeeper.AuthResponse, error) {
				return nil, expectedErr
			},
		}

		client := &GophKeeperClient{client: mockClient}
		_, err := client.Register(ctx, testReq)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestGophKeeperClient_Login(t *testing.T) {
	ctx := context.Background()
	testReq := &models.UserLoginReq{
		Username: testUser,
		Password: testPass,
	}

	t.Run("successful login", func(t *testing.T) {
		mockClient := &mockGophKeeperClient{
			loginFunc: func(ctx context.Context, in *gophkeeper.LoginRequest, opts ...grpc.CallOption) (*gophkeeper.AuthResponse, error) {
				assert.Equal(t, testUser, in.Username)
				assert.Equal(t, testPass, in.Password)
				return &gophkeeper.AuthResponse{
					UserId: testUserID,
					Token:  testToken,
					Salt:   testSalt,
				}, nil
			},
		}

		client := &GophKeeperClient{client: mockClient}
		user, err := client.Login(ctx, testReq)

		require.NoError(t, err)
		assert.Equal(t, models.UserID(testUserID), user.ID)
		assert.Equal(t, testToken, user.JWT)
		assert.Equal(t, testSalt, user.Salt)
	})

	t.Run("login error", func(t *testing.T) {
		expectedErr := errors.New("login failed")
		mockClient := &mockGophKeeperClient{
			loginFunc: func(ctx context.Context, in *gophkeeper.LoginRequest, opts ...grpc.CallOption) (*gophkeeper.AuthResponse, error) {
				return nil, expectedErr
			},
		}

		client := &GophKeeperClient{client: mockClient}
		_, err := client.Login(ctx, testReq)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestGophKeeperClient_Sync(t *testing.T) {
	ctx := context.Background()
	testTime := time.Now()
	testItems := []models.Item{
		{
			ID:        models.ItemID(testItemID),
			UserID:    models.UserID(testUserID),
			ItemType:  models.TypePassword,
			Name:      "test item",
			Metadata:  testMetadata,
			Data:      []byte("test data"),
			UpdatedAt: testTime,
			IsDeleted: false,
		},
	}

	t.Run("successful sync", func(t *testing.T) {
		mockClient := &mockGophKeeperClient{
			syncFunc: func(ctx context.Context, in *gophkeeper.SyncRequest, opts ...grpc.CallOption) (*gophkeeper.SyncResponse, error) {
				// Проверяем метаданные авторизации
				md, ok := metadata.FromOutgoingContext(ctx)
				require.True(t, ok)
				assert.Equal(t, []string{"Bearer " + testToken}, md.Get("authorization"))

				// Проверяем переданные элементы
				require.Len(t, in.Items, 1)
				assert.Equal(t, testItemID, in.Items[0].Id)
				assert.Equal(t, testUserID, in.Items[0].UserId)
				assert.Equal(t, string(models.TypePassword), in.Items[0].Type)
				assert.Equal(t, "test item", in.Items[0].Name)
				assert.Equal(t, testMetadata, in.Items[0].Metadata)
				assert.Equal(t, []byte("test data"), in.Items[0].Data)
				assert.False(t, in.Items[0].IsDeleted)

				// Возвращаем тестовые данные
				return &gophkeeper.SyncResponse{
					Items: []*gophkeeper.Item{
						{
							Id:        testItemID,
							UserId:    testUserID,
							Type:      string(models.TypePassword),
							Name:      "server item",
							Metadata:  testMetadata,
							Data:      []byte("server data"),
							UpdatedAt: timestamppb.New(testTime.Add(time.Hour)),
							IsDeleted: false,
						},
					},
				}, nil
			},
		}

		client := &GophKeeperClient{client: mockClient}
		items, err := client.Sync(ctx, testItems, testToken)

		require.NoError(t, err)
		require.Len(t, items, 1)
		assert.Equal(t, models.ItemID(testItemID), items[0].ID)
		assert.Equal(t, models.UserID(testUserID), items[0].UserID)
		assert.Equal(t, models.TypePassword, items[0].ItemType)
		assert.Equal(t, "server item", items[0].Name)
		assert.Equal(t, testMetadata, items[0].Metadata)
		assert.Equal(t, []byte("server data"), items[0].Data)
		assert.False(t, items[0].IsDeleted)
	})

	t.Run("sync error", func(t *testing.T) {
		expectedErr := errors.New("sync failed")
		mockClient := &mockGophKeeperClient{
			syncFunc: func(ctx context.Context, in *gophkeeper.SyncRequest, opts ...grpc.CallOption) (*gophkeeper.SyncResponse, error) {
				return nil, expectedErr
			},
		}

		client := &GophKeeperClient{client: mockClient}
		_, err := client.Sync(ctx, testItems, testToken)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("empty items", func(t *testing.T) {
		mockClient := &mockGophKeeperClient{
			syncFunc: func(ctx context.Context, in *gophkeeper.SyncRequest, opts ...grpc.CallOption) (*gophkeeper.SyncResponse, error) {
				assert.Empty(t, in.Items)
				return &gophkeeper.SyncResponse{Items: nil}, nil
			},
		}

		client := &GophKeeperClient{client: mockClient}
		items, err := client.Sync(ctx, nil, testToken)

		require.NoError(t, err)
		assert.Empty(t, items)
	})
}
