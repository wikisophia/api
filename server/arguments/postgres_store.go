package arguments

import (
	"database/sql"
	"log"

	"github.com/hashicorp/go-multierror"
)

// NewPostgresStore returns a Store which is used to save and load Arguments.
// The db should point to a Postgres database.
// The returned Store.Close() function will *not* close this connection, since we did not open it.
func NewPostgresStore(db *sql.DB) *PostgresStore {
	if db == nil {
		log.Fatal("A database connection is required to make an arguments.PostgresStore.")
	}
	return &PostgresStore{
		db:                           db,
		deleteStatement:              mustPrepareQuery(db, deleteQuery),
		fetchStatement:               mustPrepareQuery(db, fetchQuery),
		fetchLiveStatement:           mustPrepareQuery(db, fetchLiveQuery),
		newArgumentVersionStatement:  mustPrepareQuery(db, newArgumentVersionQuery),
		saveArgumentStatement:        mustPrepareQuery(db, saveArgumentQuery),
		saveArgumentVersionStatement: mustPrepareQuery(db, saveArgumentVersionQuery),
		saveClaimStatement:           mustPrepareQuery(db, saveClaimQuery),
		savePremiseStatement:         mustPrepareQuery(db, savePremiseQuery),
	}
}

// PostgresStore expects that {projectRoot}/postgres/scripts/create.sql
// has already been run on your database so that the expected schema exists.
type PostgresStore struct {
	db                           *sql.DB
	deleteStatement              *sql.Stmt
	fetchStatement               *sql.Stmt
	fetchLiveStatement           *sql.Stmt
	newArgumentVersionStatement  *sql.Stmt
	saveClaimStatement           *sql.Stmt
	saveArgumentStatement        *sql.Stmt
	saveArgumentVersionStatement *sql.Stmt
	savePremiseStatement         *sql.Stmt
}

// Close closes all the prepared statements used to make queries.
// It does not shut down the database connection which was passed
// into NewPostgresStore().
func (store *PostgresStore) Close() error {
	var result *multierror.Error
	result = mayAppendError(result, store.deleteStatement.Close())
	result = mayAppendError(result, store.fetchStatement.Close())
	result = mayAppendError(result, store.fetchLiveStatement.Close())
	result = mayAppendError(result, store.newArgumentVersionStatement.Close())
	result = mayAppendError(result, store.saveArgumentStatement.Close())
	result = mayAppendError(result, store.saveArgumentVersionStatement.Close())
	result = mayAppendError(result, store.saveClaimStatement.Close())
	result = mayAppendError(result, store.savePremiseStatement.Close())
	return result.ErrorOrNil()
}

func mayAppendError(existing *multierror.Error, errOrNil error) *multierror.Error {
	if errOrNil != nil {
		return multierror.Append(existing, errOrNil)
	}
	return existing
}

func mustPrepareQuery(db *sql.DB, query string) *sql.Stmt {
	statement, err := db.Prepare(query)
	if err != nil {
		log.Fatalf("Failed to prepare statement with query %s. Error was %v", query, err)
	}
	return statement
}
