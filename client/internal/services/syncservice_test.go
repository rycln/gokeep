package services

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rycln/gokeep/client/internal/services/mocks"
	"github.com/rycln/gokeep/shared/models"
	"github.com/stretchr/testify/assert"
)

func TestSyncService_SyncUserItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testUser := &models.User{
		ID:  "user123",
		JWT: "token123",
	}
	testClientItems := []models.Item{
		{ID: "item1", UserID: "user123"},
		{ID: "item2", UserID: "user123"},
	}
	testServerItems := []models.Item{
		{ID: "item1", UserID: "user123"},
		{ID: "item3", UserID: "user123"},
	}

	t.Run("successful synchronization", func(t *testing.T) {
		mockSync := mocks.NewMocksyncAPI(ctrl)
		mockStorage := mocks.NewMocksyncStorage(ctrl)

		svc := NewSyncService(mockSync, mockStorage)

		mockStorage.EXPECT().
			GetAllUserItems(gomock.Any(), testUser.ID).
			Return(testClientItems, nil)

		mockSync.EXPECT().
			Sync(gomock.Any(), testClientItems, testUser.JWT).
			Return(testServerItems, nil)

		mockStorage.EXPECT().
			ReplaceAllUserItems(gomock.Any(), testUser.ID, testServerItems).
			Return(nil)

		err := svc.SyncUserItems(context.Background(), testUser)

		assert.NoError(t, err)
	})

	t.Run("error getting local items", func(t *testing.T) {
		mockSync := mocks.NewMocksyncAPI(ctrl)
		mockStorage := mocks.NewMocksyncStorage(ctrl)

		svc := NewSyncService(mockSync, mockStorage)

		expectedErr := errors.New("storage error")

		mockStorage.EXPECT().
			GetAllUserItems(gomock.Any(), testUser.ID).
			Return(nil, expectedErr)

		err := svc.SyncUserItems(context.Background(), testUser)

		assert.EqualError(t, err, expectedErr.Error())
	})

	t.Run("error during sync with server", func(t *testing.T) {
		mockSync := mocks.NewMocksyncAPI(ctrl)
		mockStorage := mocks.NewMocksyncStorage(ctrl)

		svc := NewSyncService(mockSync, mockStorage)

		expectedErr := errors.New("sync error")

		mockStorage.EXPECT().
			GetAllUserItems(gomock.Any(), testUser.ID).
			Return(testClientItems, nil)

		mockSync.EXPECT().
			Sync(gomock.Any(), testClientItems, testUser.JWT).
			Return(nil, expectedErr)

		err := svc.SyncUserItems(context.Background(), testUser)
		assert.EqualError(t, err, expectedErr.Error())
	})

	t.Run("error saving synced items", func(t *testing.T) {
		mockSync := mocks.NewMocksyncAPI(ctrl)
		mockStorage := mocks.NewMocksyncStorage(ctrl)

		svc := NewSyncService(mockSync, mockStorage)

		expectedErr := errors.New("save error")

		mockStorage.EXPECT().
			GetAllUserItems(gomock.Any(), testUser.ID).
			Return(testClientItems, nil)

		mockSync.EXPECT().
			Sync(gomock.Any(), testClientItems, testUser.JWT).
			Return(testServerItems, nil)

		mockStorage.EXPECT().
			ReplaceAllUserItems(gomock.Any(), testUser.ID, testServerItems).
			Return(expectedErr)

		err := svc.SyncUserItems(context.Background(), testUser)
		assert.EqualError(t, err, expectedErr.Error())
	})
}
