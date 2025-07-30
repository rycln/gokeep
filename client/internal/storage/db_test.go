package storage

import (
	"context"
	"errors"
	"os"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDB(t *testing.T) {
	t.Run("successful database creation", func(t *testing.T) {
		tmpfile, err := os.CreateTemp("", "testdb.*.sqlite")
		require.NoError(t, err)
		defer os.Remove(tmpfile.Name())

		db, err := NewDB(tmpfile.Name())
		require.NoError(t, err)
		assert.NotNil(t, db)

		assert.NoError(t, db.Ping())
		assert.NoError(t, db.Close())
	})
}

func TestInitDB(t *testing.T) {
	t.Run("successful table creation", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		expectedQuery := regexp.QuoteMeta(sqlCreateItemsTable)

		mock.ExpectExec(expectedQuery).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err = InitDB(context.Background(), db)
		assert.NoError(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("table creation error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		expectedQuery := regexp.QuoteMeta(sqlCreateItemsTable)

		expectedErr := errors.New("table creation failed")
		mock.ExpectExec(expectedQuery).
			WillReturnError(expectedErr)

		err = InitDB(context.Background(), db)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("context cancellation", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		expectedQuery := regexp.QuoteMeta(sqlCreateItemsTable)

		mock.ExpectExec(expectedQuery).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err = InitDB(ctx, db)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "context canceled")
	})
}
