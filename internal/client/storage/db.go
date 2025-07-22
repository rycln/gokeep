package storage

import (
	"context"
	"database/sql"

	_ "modernc.org/sqlite"
)

func NewDB(path string) (*sql.DB, error) {
	database, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	return database, nil
}

func InitDB(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, sqlCreateItemsTable)
	if err != nil {
		return err
	}
	return nil
}
