package database

import (
	"context"

	"github.com/n25a/eavesdropper/internal/config"

	"github.com/jackc/pgx/v4/pgxpool"
)

type timescaleDB struct {
	dbPool *pgxpool.Pool
}

// NewTimescaleDB returns a new instance of timescaleDB.
func NewTimescaleDB() Database {
	return &timescaleDB{}
}

// Connect connects to the database.
func (t *timescaleDB) Connect() error {
	var err error
	ctx := context.Background()
	t.dbPool, err = pgxpool.Connect(ctx, config.C.Database.Conf.Address)
	if err != nil {
		return err
	}
	return nil
}

// CreateTable creates all tables in the database.
func (t *timescaleDB) CreateTable() error {
	panic("Not implemented")
}

// Close closes the connection with the database.
func (t *timescaleDB) Close() error {
	t.dbPool.Close()
	return nil
}

// Insert inserts the data in the database.
func (t *timescaleDB) Insert(ctx context.Context, query string, arguments ...interface{}) error {
	_, err := t.dbPool.Exec(ctx, query, arguments...)
	if err != nil {
		return err
	}
	return nil
}
