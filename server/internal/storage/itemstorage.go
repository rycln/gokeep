package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rycln/gokeep/shared/models"
)

// ItemStorage handles database operations for items.
type ItemStorage struct {
	db *sql.DB
}

// NewItemStorage creates a new ItemStorage instance.
func NewItemStorage(db *sql.DB) *ItemStorage {
	return &ItemStorage{db: db}
}

// DeleteItem removes an item from storage by ID and user ID.
func (s *ItemStorage) DeleteItem(ctx context.Context, id models.ItemID, uid models.UserID) error {
	_, err := s.db.ExecContext(ctx, sqlDeleteItem, time.Now(), id, uid)
	if err != nil {
		return err
	}
	return nil
}

// AddItem stores a new item in the database.
func (s *ItemStorage) AddItem(ctx context.Context, item *models.Item) error {
	_, err := s.db.ExecContext(
		ctx,
		sqlAddItem,
		item.ID,
		item.UserID,
		item.ItemType,
		item.Name,
		item.Metadata,
		item.Data,
		item.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

// GetUserItems retrieves all items belonging to a user.
func (s *ItemStorage) GetUserItems(ctx context.Context, uid models.UserID) (items []models.Item, err error) {
	rows, err := s.db.QueryContext(ctx, sqlGetUserItems, uid)
	if err != nil {
		return nil, err
	}
	defer func() {
		if rowsCloseErr := rows.Close(); rowsCloseErr != nil {
			err = fmt.Errorf("%v; rows close failed: %w", err, rowsCloseErr)
		}
	}()

	for rows.Next() {
		var item = models.Item{
			UserID: uid,
		}

		err = rows.Scan(
			&item.ID,
			&item.ItemType,
			&item.Name,
			&item.Metadata,
			&item.Data,
			&item.UpdatedAt,
			&item.IsDeleted,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return items, nil
}
