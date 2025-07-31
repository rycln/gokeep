package services

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rycln/gokeep/client/internal/services/mocks"
	"github.com/rycln/gokeep/shared/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testUserID = "550e8400-e29b-41d4-a716-446655440000"
	testToken  = "test.jwt.token"
	testUser   = "testuser"
	testPass   = "testpass"
	testSalt   = "salt"
)

func TestNewAuthService(t *testing.T) {
	t.Run("should create new auth service", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAPI := mocks.NewMockAuthAPI(ctrl)
		service := NewAuthService(mockAPI)

		assert.NotNil(t, service)
		assert.Equal(t, mockAPI, service.api)
	})
}

func TestUserService_UserRegister(t *testing.T) {
	ctx := context.Background()
	testReq := &models.UserRegReq{
		Username: testUser,
		Password: testPass,
		Salt:     testSalt,
	}

	expectedUser := &models.User{
		ID:   models.UserID(testUserID),
		JWT:  testToken,
		Salt: testSalt,
	}

	t.Run("successful registration", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAPI := mocks.NewMockAuthAPI(ctrl)
		service := NewAuthService(mockAPI)

		mockAPI.EXPECT().
			Register(ctx, testReq).
			Return(expectedUser, nil)

		user, err := service.UserRegister(ctx, testReq)
		require.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("registration error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAPI := mocks.NewMockAuthAPI(ctrl)
		service := NewAuthService(mockAPI)

		expectedErr := errors.New("registration failed")
		mockAPI.EXPECT().
			Register(ctx, testReq).
			Return(nil, expectedErr)

		_, err := service.UserRegister(ctx, testReq)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestUserService_UserLogin(t *testing.T) {
	ctx := context.Background()
	testReq := &models.UserLoginReq{
		Username: testUser,
		Password: testPass,
	}

	expectedUser := &models.User{
		ID:   models.UserID(testUserID),
		JWT:  testToken,
		Salt: testSalt,
	}

	t.Run("successful login", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAPI := mocks.NewMockAuthAPI(ctrl)
		service := NewAuthService(mockAPI)

		mockAPI.EXPECT().
			Login(ctx, testReq).
			Return(expectedUser, nil)

		user, err := service.UserLogin(ctx, testReq)
		require.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("login error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAPI := mocks.NewMockAuthAPI(ctrl)
		service := NewAuthService(mockAPI)

		expectedErr := errors.New("login failed")
		mockAPI.EXPECT().
			Login(ctx, testReq).
			Return(nil, expectedErr)

		_, err := service.UserLogin(ctx, testReq)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}
