package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/rycln/gokeep/client/internal/services/mocks"
	"github.com/rycln/gokeep/shared/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewItemService(t *testing.T) {
	t.Run("should create new item service", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStorage := mocks.NewMockitemStorage(ctrl)
		service := NewItemService(mockStorage)

		assert.NotNil(t, service)
		assert.Equal(t, mockStorage, service.storage)
	})
}

func TestItemService_Add(t *testing.T) {
	ctx := context.Background()
	content := []byte("test content")
	info := &models.ItemInfo{
		UserID:   "user123",
		ItemType: models.TypePassword,
	}

	t.Run("successful add", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStorage := mocks.NewMockitemStorage(ctrl)
		service := NewItemService(mockStorage)

		mockStorage.EXPECT().
			Add(ctx, gomock.Any(), content).
			Do(func(ctx context.Context, actualInfo *models.ItemInfo, actualContent []byte) {
				assert.NotEmpty(t, actualInfo.ID)
				assert.Equal(t, info.UserID, actualInfo.UserID)
				assert.Equal(t, info.ItemType, actualInfo.ItemType)
				assert.WithinDuration(t, time.Now(), actualInfo.UpdatedAt, time.Second)
			}).
			Return(nil)

		err := service.Add(ctx, info, content)
		assert.NoError(t, err)
	})

	t.Run("storage error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStorage := mocks.NewMockitemStorage(ctrl)
		service := NewItemService(mockStorage)

		expectedErr := errors.New("storage error")
		mockStorage.EXPECT().
			Add(ctx, gomock.Any(), content).
			Return(expectedErr)

		err := service.Add(ctx, info, content)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestItemService_List(t *testing.T) {
	ctx := context.Background()
	userID := models.UserID("user123")

	t.Run("successful list", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStorage := mocks.NewMockitemStorage(ctrl)
		service := NewItemService(mockStorage)

		expectedItems := []models.ItemInfo{
			{ID: "item1", UserID: userID},
			{ID: "item2", UserID: userID},
		}

		mockStorage.EXPECT().
			ListByUser(ctx, userID).
			Return(expectedItems, nil)

		items, err := service.List(ctx, userID)
		require.NoError(t, err)
		assert.Equal(t, expectedItems, items)
	})

	t.Run("storage error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStorage := mocks.NewMockitemStorage(ctrl)
		service := NewItemService(mockStorage)

		expectedErr := errors.New("storage error")
		mockStorage.EXPECT().
			ListByUser(ctx, userID).
			Return(nil, expectedErr)

		_, err := service.List(ctx, userID)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestItemService_GetContent(t *testing.T) {
	ctx := context.Background()
	itemID := models.ItemID("item123")

	t.Run("successful get content", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStorage := mocks.NewMockitemStorage(ctrl)
		service := NewItemService(mockStorage)

		expectedContent := []byte("encrypted content")
		mockStorage.EXPECT().
			GetContent(ctx, itemID).
			Return(expectedContent, nil)

		content, err := service.GetContent(ctx, itemID)
		require.NoError(t, err)
		assert.Equal(t, expectedContent, content)
	})

	t.Run("storage error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStorage := mocks.NewMockitemStorage(ctrl)
		service := NewItemService(mockStorage)

		expectedErr := errors.New("storage error")
		mockStorage.EXPECT().
			GetContent(ctx, itemID).
			Return(nil, expectedErr)

		_, err := service.GetContent(ctx, itemID)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestItemService_Delete(t *testing.T) {
	ctx := context.Background()
	itemID := models.ItemID("item123")

	t.Run("successful delete", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStorage := mocks.NewMockitemStorage(ctrl)
		service := NewItemService(mockStorage)

		mockStorage.EXPECT().
			DeleteItem(ctx, itemID).
			Return(nil)

		err := service.Delete(ctx, itemID)
		assert.NoError(t, err)
	})

	t.Run("storage error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStorage := mocks.NewMockitemStorage(ctrl)
		service := NewItemService(mockStorage)

		expectedErr := errors.New("storage error")
		mockStorage.EXPECT().
			DeleteItem(ctx, itemID).
			Return(expectedErr)

		err := service.Delete(ctx, itemID)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestItemService_Update(t *testing.T) {
	ctx := context.Background()
	content := []byte("updated content")
	info := &models.ItemInfo{
		ID:       "item123",
		UserID:   "user123",
		ItemType: models.TypePassword,
	}

	t.Run("successful update", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStorage := mocks.NewMockitemStorage(ctrl)
		service := NewItemService(mockStorage)

		mockStorage.EXPECT().
			UpdateItem(ctx, gomock.Any(), content).
			Do(func(ctx context.Context, actualInfo *models.ItemInfo, actualContent []byte) {
				assert.Equal(t, info.ID, actualInfo.ID)
				assert.Equal(t, info.UserID, actualInfo.UserID)
				assert.Equal(t, info.ItemType, actualInfo.ItemType)
				assert.WithinDuration(t, time.Now(), actualInfo.UpdatedAt, time.Second)
			}).
			Return(nil)

		err := service.Update(ctx, info, content)
		assert.NoError(t, err)
	})

	t.Run("storage error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStorage := mocks.NewMockitemStorage(ctrl)
		service := NewItemService(mockStorage)

		expectedErr := errors.New("storage error")
		mockStorage.EXPECT().
			UpdateItem(ctx, gomock.Any(), content).
			Return(expectedErr)

		err := service.Update(ctx, info, content)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}
