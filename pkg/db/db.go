// Package db provides functionality for initializing a PostgreSQL database connection.
//
// It wraps sql.Open and sql.Ping into a reusable helper for consistent database setup.
package db

import (
	"database/sql"
	"fmt"
)

// InitDB initializes and verifies a database connection using the given DSN and driver name.
//
// It returns:
//   - a valid *sql.DB instance on success;
//   - an error if connection or ping fails.
func InitDB(dsn string, driverName string) (*sql.DB, error) {
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open error: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("db.Ping error: %w", err)
	}

	return db, nil
}
