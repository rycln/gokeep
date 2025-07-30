package storage

import (
	"context"
	"database/sql"

	_ "modernc.org/sqlite" // SQLite driver
)

// NewDB creates and opens a new SQLite database connection
// path specifies the file path for the SQLite database
func NewDB(path string) (*sql.DB, error) {
	database, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	return database, nil
}

// InitDB initializes the database schema
// Creates required tables if they don't exist
func InitDB(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, sqlCreateItemsTable)
	if err != nil {
		return err
	}
	return nil
}
