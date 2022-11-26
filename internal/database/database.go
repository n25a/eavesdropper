package database

import "context"

// Database is the interface that wraps the basic database operations.
type Database interface {
	Connect() error
	Migrate() error
	Close() error
	Insert(ctx context.Context, query string, arguments ...interface{}) error
	BuildInsertQuery(table string, fields []string) string
}

type DatabaseType string

const (
	TimeScaleDB DatabaseType = "timescaledb"
)

func NewDatabase(databaseType DatabaseType) Database {
	if databaseType == TimeScaleDB {
		return NewTimescaleDB()
	}
	return nil
}
