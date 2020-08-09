package accounts

import (
	"database/sql"
	"log"
	"strings"
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
	var errs []error
	if len(errs) == 0 {
		return nil
	}
	return closeErrors(errs)
}

type closeErrors []error

func (errs closeErrors) Error() string {
	if len(errs) == 0 {
		return ""
	}

	sb := strings.Builder{}
	sb.WriteString("error(s) occurred while shutting down the accounts.Store:\n")
	for i := 0; i < len(errs); i++ {
		sb.WriteString("  ")
		sb.WriteString(errs[i].Error())
		sb.WriteString("\n")
	}
	return sb.String()
}

func mayAppendError(f func() error, errs []error) []error {
	if err := f(); err != nil {
		return append(errs, err)
	}
	return errs
}

func mustPrepareQuery(db *sql.DB, query string) *sql.Stmt {
	statement, err := db.Prepare(query)
	if err != nil {
		log.Fatalf("Failed to prepare statement with query %s. Error was %v", query, err)
	}
	return statement
}
