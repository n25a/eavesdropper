package database

// Database is the interface that wraps the basic database operations.
type Database interface {
	Connect() error
	CreateTable() error
	Close() error
	Insert(query string, data interface{}) error
}

// MapSubjectToQuery is a map that contains the subject as key and the query as value.
type MapSubjectToQuery map[string]string

type DatabaseType string

const (
	TimeScaleDB DatabaseType = "timescaledb"
)
