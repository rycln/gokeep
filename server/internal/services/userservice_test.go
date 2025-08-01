package services

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/rycln/gokeep/server/internal/contextkeys"
	"github.com/rycln/gokeep/server/internal/services/mocks"
	"github.com/rycln/gokeep/shared/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testUserID       = models.UserID("550e8400-e29b-41d4-a716-446655440000")
	testJWTToken     = "test.jwt.token"
	testPassword     = "secret"
	testPasswordHash = "hashed_secret"
	testSalt         = "salt"
)

var (
	errTest = errors.New("test error")
)

func TestUserService_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mStrg := mocks.NewMockuserStorage(ctrl)
	mHasher := mocks.NewMockpassHasher(ctrl)
	mJWT := mocks.NewMockjwtCreator(ctrl)

	t.Run("successful user creation", func(t *testing.T) {
		req := &models.UserRegReq{
			Username: "testuser",
			Password: testPassword,
			Salt:     testSalt,
		}

		gomock.InOrder(
			mHasher.EXPECT().Hash(req.Password).Return(testPasswordHash, nil),
			mStrg.EXPECT().AddUser(gomock.Any(), gomock.Any()).DoAndReturn(
				func(ctx context.Context, userDB *models.UserDB) error {
					_, err := uuid.Parse(string(userDB.ID))
					assert.NoError(t, err)
					assert.Equal(t, req.Username, userDB.Username)
					assert.Equal(t, testPasswordHash, userDB.PassHash)
					return nil
				}),
			mJWT.EXPECT().NewJWTString(gomock.Any()).DoAndReturn(
				func(userID models.UserID) (string, error) {
					_, err := uuid.Parse(string(userID))
					assert.NoError(t, err)
					return testJWTToken, nil
				}),
		)

		s := NewUserService(mStrg, mHasher, mJWT)
		user, err := s.CreateUser(context.Background(), req)
		assert.NoError(t, err)

		_, err = uuid.Parse(string(user.ID))
		assert.NoError(t, err)
		assert.Equal(t, testJWTToken, user.JWT)
	})

	t.Run("password hashing failed", func(t *testing.T) {
		req := &models.UserRegReq{
			Username: "testuser",
			Password: testPassword,
			Salt:     testSalt,
		}

		mHasher.EXPECT().Hash(req.Password).Return("", errTest)

		s := NewUserService(mStrg, mHasher, mJWT)
		_, err := s.CreateUser(context.Background(), req)
		assert.Error(t, err)
	})

	t.Run("storage add user failed", func(t *testing.T) {
		req := &models.UserRegReq{
			Username: "testuser",
			Password: testPassword,
			Salt:     testSalt,
		}

		gomock.InOrder(
			mHasher.EXPECT().Hash(req.Password).Return(testPasswordHash, nil),
			mStrg.EXPECT().AddUser(gomock.Any(), gomock.Any()).Return(errTest),
		)

		s := NewUserService(mStrg, mHasher, mJWT)
		_, err := s.CreateUser(context.Background(), req)
		assert.Error(t, err)
	})

	t.Run("JWT generation failed", func(t *testing.T) {
		req := &models.UserRegReq{
			Username: "testuser",
			Password: testPassword,
			Salt:     testSalt,
		}

		gomock.InOrder(
			mHasher.EXPECT().Hash(req.Password).Return(testPasswordHash, nil),
			mStrg.EXPECT().AddUser(gomock.Any(), gomock.Any()).Return(nil),
			mJWT.EXPECT().NewJWTString(gomock.Any()).Return("", errTest),
		)

		s := NewUserService(mStrg, mHasher, mJWT)
		_, err := s.CreateUser(context.Background(), req)
		assert.Error(t, err)
	})
}

func TestUserService_AuthUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mStrg := mocks.NewMockuserStorage(ctrl)
	mHasher := mocks.NewMockpassHasher(ctrl)
	mJWT := mocks.NewMockjwtCreator(ctrl)

	t.Run("successful authentication", func(t *testing.T) {
		req := &models.UserLoginReq{
			Username: "testuser",
			Password: testPassword,
		}

		userDB := &models.UserDB{
			ID:       models.UserID(testUserID),
			Username: req.Username,
			PassHash: testPasswordHash,
			Salt:     testSalt,
		}

		expectedUser := &models.User{
			ID:   models.UserID(testUserID),
			JWT:  testJWTToken,
			Salt: testSalt,
		}

		gomock.InOrder(
			mStrg.EXPECT().GetUserByUsername(gomock.Any(), req.Username).Return(userDB, nil),
			mHasher.EXPECT().Compare(userDB.PassHash, req.Password).Return(nil),
			mJWT.EXPECT().NewJWTString(userDB.ID).Return(testJWTToken, nil),
		)

		s := NewUserService(mStrg, mHasher, mJWT)
		user, err := s.AuthUser(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("user not found", func(t *testing.T) {
		req := &models.UserLoginReq{
			Username: "nonexistent",
			Password: testPassword,
		}

		mStrg.EXPECT().GetUserByUsername(gomock.Any(), req.Username).Return(nil, errTest)

		s := NewUserService(mStrg, mHasher, mJWT)
		_, err := s.AuthUser(context.Background(), req)
		assert.Error(t, err)
	})

	t.Run("password mismatch", func(t *testing.T) {
		req := &models.UserLoginReq{
			Username: "testuser",
			Password: "wrong_password",
		}

		userDB := &models.UserDB{
			ID:       models.UserID(testUserID),
			Username: req.Username,
			PassHash: testPasswordHash,
			Salt:     testSalt,
		}

		gomock.InOrder(
			mStrg.EXPECT().GetUserByUsername(gomock.Any(), req.Username).Return(userDB, nil),
			mHasher.EXPECT().Compare(userDB.PassHash, req.Password).Return(errTest),
		)

		s := NewUserService(mStrg, mHasher, mJWT)
		_, err := s.AuthUser(context.Background(), req)
		assert.Error(t, err)
	})

	t.Run("JWT generation failed", func(t *testing.T) {
		req := &models.UserLoginReq{
			Username: "testuser",
			Password: testPassword,
		}

		userDB := &models.UserDB{
			ID:       models.UserID(testUserID),
			Username: req.Username,
			PassHash: testPasswordHash,
			Salt:     testSalt,
		}

		gomock.InOrder(
			mStrg.EXPECT().GetUserByUsername(gomock.Any(), req.Username).Return(userDB, nil),
			mHasher.EXPECT().Compare(userDB.PassHash, req.Password).Return(nil),
			mJWT.EXPECT().NewJWTString(userDB.ID).Return("", errTest),
		)

		s := NewUserService(mStrg, mHasher, mJWT)
		_, err := s.AuthUser(context.Background(), req)
		assert.Error(t, err)
	})
}

func TestUserService_GetUserIDFromCtx(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := &UserService{}

	testUserID := models.UserID("550e8400-e29b-41d4-a716-446655440000")

	t.Run("successfully get user id from context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), contextkeys.UserID, testUserID)

		uid, err := service.GetUserIDFromCtx(ctx)

		require.NoError(t, err)
		assert.Equal(t, testUserID, uid)
	})

	t.Run("error when no user id in context", func(t *testing.T) {
		ctx := context.Background()

		uid, err := service.GetUserIDFromCtx(ctx)

		require.Error(t, err)
		assert.Equal(t, models.UserID(""), uid)
		assert.Equal(t, errNoUserID, err)
	})

	t.Run("error when wrong value type in context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), contextkeys.UserID, "wrong-type")

		uid, err := service.GetUserIDFromCtx(ctx)

		require.Error(t, err)
		assert.Equal(t, models.UserID(""), uid)
		assert.Equal(t, errNoUserID, err)
	})
}
