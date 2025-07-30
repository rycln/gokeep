package storage

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Connection pool configuration constants
const (
	maxOpenConns    = 0               // Unlimited open connections (0 means no limit)
	maxIdleConns    = 10              // Maximum idle connections in pool
	maxIdleTime     = 3 * time.Minute // Maximum time a connection can remain idle
	maxConnLifetime = 0               // Maximum connection lifetime (0 means no limit)
)

// NewDB creates and configures a new PostgreSQL database connection pool
func NewDB(uri string) (*sql.DB, error) {
	database, err := sql.Open("pgx", uri)
	if err != nil {
		return nil, err
	}

	database.SetMaxOpenConns(maxOpenConns)
	database.SetMaxIdleConns(maxIdleConns)
	database.SetConnMaxIdleTime(maxIdleTime)
	database.SetConnMaxLifetime(maxConnLifetime)

	if err := database.Ping(); err != nil {
		return nil, err
	}

	return database, nil
}
