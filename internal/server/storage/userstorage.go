package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rycln/gokeep/internal/shared/models"
)

type UserStorage struct {
	db *sql.DB
}

func NewUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{db: db}
}

func (s *UserStorage) AddUser(ctx context.Context, user *models.UserDB) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			err = fmt.Errorf("%v; rollback failed: %w", err, rollbackErr)
		}
	}()

	_, err = tx.ExecContext(ctx, sqlAddUser, user.ID, user.Username, user.PassHash)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return newErrUsernameConflict(ErrUsernameConflict)
		}
		return err
	}

	return tx.Commit()
}

func (s *UserStorage) GetUserByUsername(ctx context.Context, username string) (*models.UserDB, error) {
	row := s.db.QueryRowContext(ctx, sqlGetUserByUsername, username)
	var userDB models.UserDB
	err := row.Scan(&userDB.ID, &userDB.Username, &userDB.PassHash)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, newErrNoUser(ErrNoUser)
	}
	if err != nil {
		return nil, err
	}
	return &userDB, nil
}
