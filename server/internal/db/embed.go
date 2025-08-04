// Package db contains database schema migration files.
package db

import (
	"embed"
)

//go:embed migrations/*.sql
var MigrationsFS embed.FS
