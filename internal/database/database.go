package database

import (
	"context"
	"errors"
)

// Database is the interface that wraps the basic database operations.
type Database interface {
	Connect() error
	Migrate() error
	Close() error
	Insert(ctx context.Context, query string, arguments ...interface{}) error
	Get(ctx context.Context, query string, dest interface{}, arguments ...interface{}) (interface{}, error)
	BuildInsertQuery(table string, fields []string) string
}

// ErrNotFound is the error returned when the database does not find the data.
var ErrNotFound = errors.New("not found")

// DatabaseType is the type of database.
type DatabaseType string

// Database types
const (
	TimeScaleDB DatabaseType = "timescaledb"
)

// NewDatabase returns a new database instance.
func NewDatabase(databaseType DatabaseType) Database {
	if databaseType == TimeScaleDB {
		return NewTimescaleDB()
	}
	return nil
}
