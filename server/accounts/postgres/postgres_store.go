package postgres

import (
	"database/sql"
	"log"

	"github.com/hashicorp/go-multierror"
)

// NewPostgresStore returns a Store which can manage accounts.
// The db should point to a Postgres database.
// The returned Store.Close() function will *not* close this connection, since we did not open it.
func NewPostgresStore(db *sql.DB) *PostgresStore {
	if db == nil {
		log.Fatal("A database connection is required to make an accounts.PostgresStore.")
	}
	return &PostgresStore{
		db: db,
	}
}

// PostgresStore saves account info in Postgres.
type PostgresStore struct {
	db *sql.DB
}

// Close closes all the prepared statements used to make queries.
// It does not shut down the database connection which was passed
// into NewPostgresStore().
func (store *PostgresStore) Close() error {
	var err *multierror.Error
	// TODO: Close prepared statements & append errors here
	return err.ErrorOrNil()
}
