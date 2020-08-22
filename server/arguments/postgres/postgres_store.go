package postgres

import (
	"database/sql"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/wikisophia/api/server/postgres"
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
		deleteStatement:              postgres.MustPrepareQuery(db, deleteQuery),
		fetchStatement:               postgres.MustPrepareQuery(db, fetchQuery),
		fetchLiveStatement:           postgres.MustPrepareQuery(db, fetchLiveQuery),
		newArgumentVersionStatement:  postgres.MustPrepareQuery(db, newArgumentVersionQuery),
		saveArgumentStatement:        postgres.MustPrepareQuery(db, saveArgumentQuery),
		saveArgumentVersionStatement: postgres.MustPrepareQuery(db, saveArgumentVersionQuery),
		saveClaimStatement:           postgres.MustPrepareQuery(db, saveClaimQuery),
		savePremiseStatement:         postgres.MustPrepareQuery(db, savePremiseQuery),
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
	result = multierror.Append(result, store.deleteStatement.Close())
	result = multierror.Append(result, store.fetchStatement.Close())
	result = multierror.Append(result, store.fetchLiveStatement.Close())
	result = multierror.Append(result, store.newArgumentVersionStatement.Close())
	result = multierror.Append(result, store.saveArgumentStatement.Close())
	result = multierror.Append(result, store.saveArgumentVersionStatement.Close())
	result = multierror.Append(result, store.saveClaimStatement.Close())
	result = multierror.Append(result, store.savePremiseStatement.Close())
	return result.ErrorOrNil()
}
