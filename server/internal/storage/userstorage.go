package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rycln/gokeep/shared/models"
)

// UserStorage manages persistence operations for user entities
type UserStorage struct {
	db *sql.DB
}

// NewUserStorage creates a new UserStorage instance with the given database connection
func NewUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{db: db}
}

// AddUser persists a new user to the database
func (s *UserStorage) AddUser(ctx context.Context, user *models.UserDB) error {
	_, err := s.db.ExecContext(ctx, sqlAddUser, user.ID, user.Username, user.PassHash, user.Salt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return newErrUsernameConflict(ErrUsernameConflict)
		}
		return err
	}
	return nil
}

// GetUserByUsername retrieves a user by their username
func (s *UserStorage) GetUserByUsername(ctx context.Context, username string) (*models.UserDB, error) {
	row := s.db.QueryRowContext(ctx, sqlGetUserByUsername, username)

	var userDB models.UserDB
	err := row.Scan(&userDB.ID, &userDB.Username, &userDB.PassHash, &userDB.Salt)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, newErrNoUser(ErrNoUser)
	case err != nil:
		return nil, err
	default:
		return &userDB, nil
	}
}
