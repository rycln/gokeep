package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/rycln/gokeep/shared/models"
)

// ItemStorage handles persistent storage operations for items
type ItemStorage struct {
	db *sql.DB // Database connection
}

// NewItemStorage creates a new ItemStorage instance
func NewItemStorage(db *sql.DB) *ItemStorage {
	return &ItemStorage{db: db}
}

// Add stores a new item with its content and metadata
func (s *ItemStorage) Add(ctx context.Context, info *models.ItemInfo, content []byte) error {
	_, err := s.db.ExecContext(
		ctx,
		sqlAddItem,
		info.ID,
		info.UserID,
		info.ItemType,
		info.Name,
		content,
		info.Metadata,
	)
	return err
}

// ListByUser retrieves all item metadata for a specific user
func (s *ItemStorage) ListByUser(ctx context.Context, uid models.UserID) ([]models.ItemInfo, error) {
	rows, err := s.db.QueryContext(ctx, sqlGetUserItemsInfo, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.ItemInfo
	for rows.Next() {
		var info models.ItemInfo
		if err := rows.Scan(
			&info.ID,
			&info.UserID,
			&info.ItemType,
			&info.Name,
			&info.Metadata,
			&info.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, info)
	}

	return items, rows.Err()
}

// GetContent retrieves the encrypted content of a specific item
func (s *ItemStorage) GetContent(ctx context.Context, id models.ItemID) ([]byte, error) {
	var content []byte
	err := s.db.QueryRowContext(ctx, sqlGetItemByID, id).Scan(&content)
	return content, err
}

// DeleteItem removes an item from storage by its ID
func (s *ItemStorage) DeleteItem(ctx context.Context, id models.ItemID) error {
	_, err := s.db.ExecContext(ctx, sqlDeleteItem, id)
	return err
}

// UpdateItem modifies an existing item's data and content
func (s *ItemStorage) UpdateItem(ctx context.Context, info *models.ItemInfo, content []byte) error {
	_, err := s.db.ExecContext(
		ctx,
		sqlUpdateItem,
		info.Name,
		info.Metadata,
		info.UpdatedAt,
		content,
		info.UserID,
		info.ID,
	)
	return err
}

// GetAllUserItems retrieves all items (including content and deleted status) for a specific user
func (s *ItemStorage) GetAllUserItems(ctx context.Context, uid models.UserID) ([]models.Item, error) {
	rows, err := s.db.QueryContext(ctx, sqlGetAllUserItems, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.ItemType,
			&item.Name,
			&item.Metadata,
			&item.Data,
			&item.UpdatedAt,
			&item.IsDeleted,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

// ReplaceAllUserItems completely replaces all items for a user in a single transaction
// First deletes all existing items, then inserts the new ones
func (s *ItemStorage) ReplaceAllUserItems(ctx context.Context, uid models.UserID, items []models.Item) (err error) {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			err = fmt.Errorf("%v; rollback failed: %w", err, rollbackErr)
		}
	}()

	if _, err := tx.Exec(sqlDeleteUserItems, uid); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to clear items table: %w", err)
	}

	stmt, err := tx.Prepare(sqlAddUserItems)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	defer stmt.Close()

	for _, item := range items {
		_, err := stmt.Exec(
			item.ID,
			item.UserID,
			item.ItemType,
			item.Name,
			item.Metadata,
			item.Data,
			item.UpdatedAt,
			item.IsDeleted,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert item %s: %w", item.ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
