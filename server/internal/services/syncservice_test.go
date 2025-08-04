package services

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rycln/gokeep/server/internal/services/mocks"
	"github.com/rycln/gokeep/shared/models"
	"github.com/stretchr/testify/assert"
)

func TestNewSyncService(t *testing.T) {
	t.Run("should create new SyncService instance", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStorage := mocks.NewMockitemStorage(ctrl)
		mockAuth := mocks.NewMockuidFetcher(ctrl)

		service := NewSyncService(mockStorage, mockAuth)
		assert.NotNil(t, service)
	})
}

func TestSyncItems(t *testing.T) {
	t.Run("should successfully sync items", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStorage := mocks.NewMockitemStorage(ctrl)
		mockAuth := mocks.NewMockuidFetcher(ctrl)

		userID := models.UserID("user123")
		ctx := context.Background()
		reqItems := []models.Item{
			{ID: models.ItemID("item1"), IsDeleted: false},
			{ID: models.ItemID("item2"), IsDeleted: true},
		}
		resItems := []models.Item{
			{ID: models.ItemID("item1"), Name: "Synced Item"},
		}

		mockAuth.EXPECT().
			GetUserIDFromCtx(ctx).
			Return(userID, nil)

		mockStorage.EXPECT().
			AddItem(ctx, &reqItems[0]).
			Return(nil)

		mockStorage.EXPECT().
			DeleteItem(ctx, models.ItemID("item2"), userID).
			Return(nil)

		mockStorage.EXPECT().
			GetUserItems(ctx, userID).
			Return(resItems, nil)

		service := NewSyncService(mockStorage, mockAuth)
		result, err := service.SyncItems(ctx, reqItems)

		assert.NoError(t, err)
		assert.Equal(t, resItems, result)
	})

	t.Run("should return error when failed to get user ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStorage := mocks.NewMockitemStorage(ctrl)
		mockAuth := mocks.NewMockuidFetcher(ctrl)

		testErr := errors.New("auth error")
		mockAuth.EXPECT().
			GetUserIDFromCtx(gomock.Any()).
			Return(models.UserID(""), testErr)

		service := NewSyncService(mockStorage, mockAuth)
		_, err := service.SyncItems(context.Background(), []models.Item{})

		assert.Equal(t, testErr, err)
	})

	t.Run("should return error when failed to add item", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStorage := mocks.NewMockitemStorage(ctrl)
		mockAuth := mocks.NewMockuidFetcher(ctrl)

		userID := models.UserID("user123")
		testErr := errors.New("add error")
		item := models.Item{ID: "item1", IsDeleted: false}

		mockAuth.EXPECT().
			GetUserIDFromCtx(gomock.Any()).
			Return(userID, nil)

		mockStorage.EXPECT().
			AddItem(gomock.Any(), &item).
			Return(testErr)

		service := NewSyncService(mockStorage, mockAuth)
		_, err := service.SyncItems(context.Background(), []models.Item{item})

		assert.Equal(t, testErr, err)
	})

	t.Run("should return error when failed to delete item", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStorage := mocks.NewMockitemStorage(ctrl)
		mockAuth := mocks.NewMockuidFetcher(ctrl)

		userID := models.UserID("user123")
		testErr := errors.New("delete error")
		item := models.Item{ID: models.ItemID("item1"), IsDeleted: true}

		mockAuth.EXPECT().
			GetUserIDFromCtx(gomock.Any()).
			Return(userID, nil)

		mockStorage.EXPECT().
			DeleteItem(gomock.Any(), models.ItemID("item1"), userID).
			Return(testErr)

		service := NewSyncService(mockStorage, mockAuth)
		_, err := service.SyncItems(context.Background(), []models.Item{item})

		assert.Equal(t, testErr, err)
	})

	t.Run("should return error when failed to get user items", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStorage := mocks.NewMockitemStorage(ctrl)
		mockAuth := mocks.NewMockuidFetcher(ctrl)

		userID := models.UserID("user123")
		testErr := errors.New("fetch error")
		item := models.Item{ID: "item1", IsDeleted: false}

		mockAuth.EXPECT().
			GetUserIDFromCtx(gomock.Any()).
			Return(userID, nil)

		mockStorage.EXPECT().
			AddItem(gomock.Any(), &item).
			Return(nil)

		mockStorage.EXPECT().
			GetUserItems(gomock.Any(), userID).
			Return(nil, testErr)

		service := NewSyncService(mockStorage, mockAuth)
		_, err := service.SyncItems(context.Background(), []models.Item{item})

		assert.Equal(t, testErr, err)
	})

	t.Run("should handle empty items list", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStorage := mocks.NewMockitemStorage(ctrl)
		mockAuth := mocks.NewMockuidFetcher(ctrl)

		userID := models.UserID("user123")
		resItems := []models.Item{}

		mockAuth.EXPECT().
			GetUserIDFromCtx(gomock.Any()).
			Return(userID, nil)

		mockStorage.EXPECT().
			GetUserItems(gomock.Any(), userID).
			Return(resItems, nil)

		service := NewSyncService(mockStorage, mockAuth)
		result, err := service.SyncItems(context.Background(), []models.Item{})

		assert.NoError(t, err)
		assert.Equal(t, resItems, result)
	})
}
