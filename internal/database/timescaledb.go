package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/golang-migrate/migrate"
	libmigratepostgres "github.com/golang-migrate/migrate/database/postgres"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/n25a/eavesdropper/internal/config"
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

// Close closes the connection with the database.
func (t *timescaleDB) Close() error {
	t.dbPool.Close()
	return nil
}

// CreateTable creates all tables in the database.
func (t *timescaleDB) CreateTable() error {
	// TODO: Create migration files from schema and migrate them.
	if config.C.Database.MigrationPath == "" {
		return errors.New("migration path is empty")
	}

	// TODO: create data source name from config.
	db, err := sql.Open("postgres", "postgres://localhost:5432/database?sslmode=enable")
	if err != nil {
		return err
	}

	driver, err := libmigratepostgres.WithInstance(db, &libmigratepostgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:///"+config.C.Database.MigrationPath,
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}

	err = m.Up()

	return err
}

// Insert inserts the data in the database.
func (t *timescaleDB) Insert(ctx context.Context, query string, arguments ...interface{}) error {
	_, err := t.dbPool.Exec(ctx, query, arguments...)
	if err != nil {
		return err
	}
	return nil
}
