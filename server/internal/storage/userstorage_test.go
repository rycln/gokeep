package storage

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rycln/gokeep/shared/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testUserID = "550e8400-e29b-41d4-a716-446655440000"
	testSalt   = "salt"
)

var (
	errTest = errors.New("test error")
)

func TestUserStorage_AddUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	strg := NewUserStorage(db)

	testUser := &models.UserDB{
		ID:       testUserID,
		Username: "testuser",
		PassHash: "hashed_password",
		Salt:     testSalt,
	}

	expectedQuery := regexp.QuoteMeta(sqlAddUser)

	t.Run("successful user creation", func(t *testing.T) {
		mock.ExpectExec(expectedQuery).
			WithArgs(testUser.ID, testUser.Username, testUser.PassHash, testUser.Salt).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := strg.AddUser(context.Background(), testUser)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("username conflict error", func(t *testing.T) {
		pgErr := &pgconn.PgError{
			Code: pgerrcode.UniqueViolation,
		}

		mock.ExpectExec(expectedQuery).
			WithArgs(testUser.ID, testUser.Username, testUser.PassHash, testUser.Salt).
			WillReturnError(pgErr)

		err := strg.AddUser(context.Background(), testUser)
		assert.ErrorIs(t, err, ErrUsernameConflict)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("general database error", func(t *testing.T) {
		mock.ExpectExec(expectedQuery).
			WithArgs(testUser.ID, testUser.Username, testUser.PassHash, testUser.Salt).
			WillReturnError(errTest)

		err := strg.AddUser(context.Background(), testUser)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserStorage_GetUserByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	strg := NewUserStorage(db)

	testUser := &models.UserDB{
		ID:       testUserID,
		Username: "testuser",
		PassHash: "hashed_password",
		Salt:     testSalt,
	}

	expectedQuery := regexp.QuoteMeta(sqlGetUserByUsername)

	t.Run("successful user retrieval", func(t *testing.T) {
		rows := mock.NewRows([]string{"id", "username", "pass_hash", "salt"}).
			AddRow(testUser.ID, testUser.Username, testUser.PassHash, testUser.Salt)

		mock.ExpectQuery(expectedQuery).
			WithArgs(testUser.Username).
			WillReturnRows(rows)

		user, err := strg.GetUserByUsername(context.Background(), testUser.Username)
		assert.NoError(t, err)
		assert.Equal(t, testUser, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("user not found error", func(t *testing.T) {
		mock.ExpectQuery(expectedQuery).
			WithArgs(testUser.Username).
			WillReturnError(sql.ErrNoRows)

		_, err := strg.GetUserByUsername(context.Background(), testUser.Username)
		assert.ErrorIs(t, err, ErrNoUser)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("general database error", func(t *testing.T) {
		mock.ExpectQuery(expectedQuery).
			WithArgs(testUser.Username).
			WillReturnError(errTest)

		_, err := strg.GetUserByUsername(context.Background(), testUser.Username)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
