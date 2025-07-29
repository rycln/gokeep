package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rycln/gokeep/internal/shared/models"
)

type ItemStorage struct {
	db *sql.DB
}

func NewItemStorage(db *sql.DB) *ItemStorage {
	return &ItemStorage{
		db: db,
	}
}

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
	if err != nil {
		return err
	}

	return nil
}

func (s *ItemStorage) ListByUser(ctx context.Context, uid models.UserID) (itemsinfo []models.ItemInfo, err error) {
	rows, err := s.db.QueryContext(ctx, sqlGetUserItemsInfo, uid)
	if err != nil {
		return nil, err
	}
	defer func() {
		if rowsCloseErr := rows.Close(); rowsCloseErr != nil {
			err = fmt.Errorf("%v; rows close failed: %w", err, rowsCloseErr)
		}
	}()

	for rows.Next() {
		var info models.ItemInfo

		err = rows.Scan(
			&info.ID,
			&info.UserID,
			&info.ItemType,
			&info.Name,
			&info.Metadata,
			&info.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		itemsinfo = append(itemsinfo, info)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return itemsinfo, nil
}

func (s *ItemStorage) GetContent(ctx context.Context, id models.ItemID) ([]byte, error) {
	row := s.db.QueryRowContext(ctx, sqlGetItemByID, id)

	var content []byte

	err := row.Scan(&content)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (s *ItemStorage) DeleteItem(ctx context.Context, id models.ItemID) error {
	_, err := s.db.ExecContext(ctx, sqlDeleteItem, id)
	if err != nil {
		return err
	}

	return nil
}
