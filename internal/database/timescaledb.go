package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
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
	t.dbPool, err = pgxpool.Connect(ctx, t.dsn())
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

// Migrate creates all tables in the database.
func (t *timescaleDB) Migrate() error {
	// TODO: Create migration files from schema and migrate them.
	if config.C.Database.MigrationPath == "" {
		return errors.New("migration path is empty")
	}

	db, err := sql.Open("postgres", t.dsn())
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

// Get returns the data from the database.
func (t *timescaleDB) Get(ctx context.Context, query string, dest interface{}, arguments ...interface{}) (interface{}, error) {
	err := t.dbPool.QueryRow(ctx, query, arguments...).Scan(dest)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return dest, nil
}

// BuildInsertQuery builder of the insert query.
func (t *timescaleDB) BuildInsertQuery(table string, fields []string) string {
	dialect := goqu.Dialect("postgres")
	ds := dialect.Insert(table)

	var records goqu.Record
	for _, f := range fields {
		records[f] = "?"
	}
	ds = ds.Rows(records)
	query, _, _ := ds.ToSQL()

	query = strings.ReplaceAll(query, "'", "")

	return query
}

// dsn return data source name
func (timescaleDB) dsn() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?parseTime=true&multiStatements=true&interpolateParams=true&collation=utf8mb4_general_ci",
		config.C.Database.Conf.User, config.C.Database.Conf.Password, config.C.Database.Conf.Host,
		config.C.Database.Conf.Port, config.C.Database.Conf.DatabaseName,
	)
}
