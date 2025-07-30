package storage

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rycln/gokeep/shared/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewItemStorage(t *testing.T) {
	t.Run("should create new item storage", func(t *testing.T) {
		db, _, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		storage := NewItemStorage(db)
		assert.NotNil(t, storage)
		assert.Equal(t, db, storage.db)
	})
}

func TestItemStorage_Add(t *testing.T) {
	ctx := context.Background()
	testInfo := &models.ItemInfo{
		ID:       "item123",
		UserID:   "user123",
		ItemType: models.TypePassword,
		Name:     "test item",
		Metadata: "metadata",
	}
	testContent := []byte("encrypted content")

	expectedQuery := regexp.QuoteMeta(sqlAddItem)

	t.Run("successful add", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		storage := NewItemStorage(db)

		mock.ExpectExec(expectedQuery).
			WithArgs(
				testInfo.ID,
				testInfo.UserID,
				testInfo.ItemType,
				testInfo.Name,
				testContent,
				testInfo.Metadata,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = storage.Add(ctx, testInfo, testContent)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		storage := NewItemStorage(db)

		expectedErr := errors.New("database error")
		mock.ExpectExec(expectedQuery).
			WithArgs(
				testInfo.ID,
				testInfo.UserID,
				testInfo.ItemType,
				testInfo.Name,
				testContent,
				testInfo.Metadata,
			).
			WillReturnError(expectedErr)

		err = storage.Add(ctx, testInfo, testContent)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestItemStorage_ListByUser(t *testing.T) {
	ctx := context.Background()
	userID := models.UserID("user123")

	expectedQuery := regexp.QuoteMeta(sqlGetUserItemsInfo)

	t.Run("successful list", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		storage := NewItemStorage(db)

		expectedItems := []models.ItemInfo{
			{
				ID:        "item1",
				UserID:    userID,
				ItemType:  models.TypePassword,
				Name:      "item 1",
				Metadata:  "metadata1",
				UpdatedAt: time.Now(),
			},
			{
				ID:        "item2",
				UserID:    userID,
				ItemType:  models.TypeCard,
				Name:      "item 2",
				Metadata:  "metadata2",
				UpdatedAt: time.Now(),
			},
		}

		rows := sqlmock.NewRows([]string{
			"id", "user_id", "type", "name", "metadata", "updated_at",
		}).
			AddRow(
				expectedItems[0].ID,
				expectedItems[0].UserID,
				expectedItems[0].ItemType,
				expectedItems[0].Name,
				expectedItems[0].Metadata,
				expectedItems[0].UpdatedAt,
			).
			AddRow(
				expectedItems[1].ID,
				expectedItems[1].UserID,
				expectedItems[1].ItemType,
				expectedItems[1].Name,
				expectedItems[1].Metadata,
				expectedItems[1].UpdatedAt,
			)

		mock.ExpectQuery(expectedQuery).
			WithArgs(userID).
			WillReturnRows(rows)

		items, err := storage.ListByUser(ctx, userID)
		require.NoError(t, err)
		assert.Equal(t, expectedItems, items)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		storage := NewItemStorage(db)

		expectedErr := errors.New("database error")
		mock.ExpectQuery(expectedQuery).
			WithArgs(userID).
			WillReturnError(expectedErr)

		_, err = storage.ListByUser(ctx, userID)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		storage := NewItemStorage(db)

		rows := sqlmock.NewRows([]string{
			"id", "user_id", "type", "name", "metadata", "updated_at",
		}).
			AddRow("item1", "user123", "invalid_type", "item 1", "{}", "{}")

		mock.ExpectQuery(expectedQuery).
			WithArgs(userID).
			WillReturnRows(rows)

		_, err = storage.ListByUser(ctx, userID)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestItemStorage_GetContent(t *testing.T) {
	ctx := context.Background()
	itemID := models.ItemID("item123")

	expectedQuery := regexp.QuoteMeta(sqlGetItemByID)

	t.Run("successful get content", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		storage := NewItemStorage(db)

		expectedContent := []byte("encrypted content")
		mock.ExpectQuery(expectedQuery).
			WithArgs(itemID).
			WillReturnRows(sqlmock.NewRows([]string{"content"}).AddRow(expectedContent))

		content, err := storage.GetContent(ctx, itemID)
		require.NoError(t, err)
		assert.Equal(t, expectedContent, content)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		storage := NewItemStorage(db)

		expectedErr := errors.New("database error")
		mock.ExpectQuery(expectedQuery).
			WithArgs(itemID).
			WillReturnError(expectedErr)

		_, err = storage.GetContent(ctx, itemID)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		storage := NewItemStorage(db)

		mock.ExpectQuery(expectedQuery).
			WithArgs(itemID).
			WillReturnError(sql.ErrNoRows)

		_, err = storage.GetContent(ctx, itemID)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestItemStorage_DeleteItem(t *testing.T) {
	ctx := context.Background()
	itemID := models.ItemID("item123")

	expectedQuery := regexp.QuoteMeta(sqlDeleteItem)

	t.Run("successful delete", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		storage := NewItemStorage(db)

		mock.ExpectExec(expectedQuery).
			WithArgs(itemID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err = storage.DeleteItem(ctx, itemID)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		storage := NewItemStorage(db)

		expectedErr := errors.New("database error")
		mock.ExpectExec(expectedQuery).
			WithArgs(itemID).
			WillReturnError(expectedErr)

		err = storage.DeleteItem(ctx, itemID)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestItemStorage_UpdateItem(t *testing.T) {
	ctx := context.Background()
	testInfo := &models.ItemInfo{
		ID:        "item123",
		UserID:    "user123",
		Name:      "updated name",
		Metadata:  "metadata",
		UpdatedAt: time.Now(),
	}
	testContent := []byte("updated content")

	expectedQuery := regexp.QuoteMeta(sqlUpdateItem)

	t.Run("successful update", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		storage := NewItemStorage(db)

		mock.ExpectExec(expectedQuery).
			WithArgs(
				testInfo.Name,
				testInfo.Metadata,
				testInfo.UpdatedAt,
				testContent,
				testInfo.UserID,
				testInfo.ID,
			).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err = storage.UpdateItem(ctx, testInfo, testContent)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		storage := NewItemStorage(db)

		expectedErr := errors.New("database error")
		mock.ExpectExec(expectedQuery).
			WithArgs(
				testInfo.Name,
				testInfo.Metadata,
				testInfo.UpdatedAt,
				testContent,
				testInfo.UserID,
				testInfo.ID,
			).
			WillReturnError(expectedErr)

		err = storage.UpdateItem(ctx, testInfo, testContent)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
