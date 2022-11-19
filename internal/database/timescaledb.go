package database

type timescaleDB struct{}

// NewTimescaleDB returns a new instance of timescaleDB.
func NewTimescaleDB() Database {
	return &timescaleDB{}
}

// Connect connects to the database.
func (t *timescaleDB) Connect() error {
	panic("Not implemented")
}

// CreateTable creates all tables in the database.
func (t *timescaleDB) CreateTable() error {
	panic("Not implemented")
}

// Close closes the connection with the database.
func (t *timescaleDB) Close() error {
	panic("Not implemented")
}

// Insert inserts the data in the database.
func (t *timescaleDB) Insert(query string, data interface{}) error {
	panic("Not implemented")
}
