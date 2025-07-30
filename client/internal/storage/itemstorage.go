package storage

import (
	"context"
	"database/sql"

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
