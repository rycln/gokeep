package storage

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rycln/gokeep/shared/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testItemID = "550e8400-e29b-41d4-a716-446655440001"
)

func TestItemStorage_AddItem(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	strg := NewItemStorage(db)

	testItem := &models.Item{
		ID:        testItemID,
		UserID:    testUserID,
		ItemType:  "note",
		Name:      "test item",
		Metadata:  "{}",
		Data:      []byte("test data"),
		UpdatedAt: time.Now(),
	}

	expectedQuery := regexp.QuoteMeta(sqlAddItem)

	t.Run("successful item creation", func(t *testing.T) {
		mock.ExpectExec(expectedQuery).
			WithArgs(
				testItem.ID,
				testItem.UserID,
				testItem.ItemType,
				testItem.Name,
				testItem.Metadata,
				testItem.Data,
				testItem.UpdatedAt,
			).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := strg.AddItem(context.Background(), testItem)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("general database error", func(t *testing.T) {
		mock.ExpectExec(expectedQuery).
			WithArgs(
				testItem.ID,
				testItem.UserID,
				testItem.ItemType,
				testItem.Name,
				testItem.Metadata,
				testItem.Data,
				testItem.UpdatedAt,
			).
			WillReturnError(errTest)

		err := strg.AddItem(context.Background(), testItem)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestItemStorage_DeleteItem(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	strg := NewItemStorage(db)

	expectedQuery := regexp.QuoteMeta(sqlDeleteItem)

	t.Run("successful item deletion", func(t *testing.T) {
		mock.ExpectExec(expectedQuery).
			WithArgs(testItemID, testUserID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := strg.DeleteItem(context.Background(), testItemID, testUserID)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("general database error", func(t *testing.T) {
		mock.ExpectExec(expectedQuery).
			WithArgs(testItemID, testUserID).
			WillReturnError(errTest)

		err := strg.DeleteItem(context.Background(), testItemID, testUserID)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestItemStorage_GetUserItems(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	strg := NewItemStorage(db)

	testTime := time.Now()
	testItems := []models.Item{
		{
			ID:        testItemID,
			UserID:    testUserID,
			ItemType:  "note",
			Name:      "test item 1",
			Metadata:  "{}",
			Data:      []byte("test data 1"),
			UpdatedAt: testTime,
			IsDeleted: false,
		},
		{
			ID:        "550e8400-e29b-41d4-a716-446655440002",
			UserID:    testUserID,
			ItemType:  "password",
			Name:      "test item 2",
			Metadata:  "{}",
			Data:      []byte("test data 2"),
			UpdatedAt: testTime.Add(time.Hour),
			IsDeleted: false,
		},
	}

	expectedQuery := regexp.QuoteMeta(sqlGetUserItems)

	t.Run("successful items retrieval", func(t *testing.T) {
		rows := mock.NewRows([]string{
			"id", "item_type", "name", "metadata", "data", "updated_at", "is_deleted",
		})
		for _, item := range testItems {
			rows.AddRow(
				item.ID,
				item.ItemType,
				item.Name,
				item.Metadata,
				item.Data,
				item.UpdatedAt,
				item.IsDeleted,
			)
		}

		mock.ExpectQuery(expectedQuery).
			WithArgs(testUserID).
			WillReturnRows(rows)

		items, err := strg.GetUserItems(context.Background(), testUserID)
		assert.NoError(t, err)
		assert.Equal(t, testItems, items)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("general database error", func(t *testing.T) {
		mock.ExpectQuery(expectedQuery).
			WithArgs(testUserID).
			WillReturnError(errTest)

		_, err := strg.GetUserItems(context.Background(), testUserID)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
