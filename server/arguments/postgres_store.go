package arguments

import (
	"database/sql"
	"log"
	"strings"
)

// NewPostgresStore returns a Store which is used to save and load Arguments.
// The db should point to a Postgres database.
// The returned Store.Close() function will *not* close this connection, since we did not open it.
func NewPostgresStore(db *sql.DB) *Store {
	if db == nil {
		log.Fatalf("A database connection is required to make a Store.")
	}
	return &Store{
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

// The Store expects that {projectRoot}/postgres/scripts/create.sql
// has already been run on your database so that the expected schema exists.
type Store struct {
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

type closeErrors []error

func (errs closeErrors) Error() string {
	if len(errs) == 0 {
		return ""
	}

	sb := strings.Builder{}
	sb.WriteString("error(s) occurred while shutting down the postgres.Store:\n")
	for i := 0; i < len(errs); i++ {
		sb.WriteString("  ")
		sb.WriteString(errs[i].Error())
		sb.WriteString("\n")
	}
	return sb.String()
}

// Close closes all the prepared statements used to make queries.
// It does not shut down the database connection which was passed
// into NewPostgresStore().
func (store *Store) Close() error {
	var errs []error
	errs = mayAppendError(store.deleteStatement.Close, errs)
	errs = mayAppendError(store.fetchStatement.Close, errs)
	errs = mayAppendError(store.fetchLiveStatement.Close, errs)
	errs = mayAppendError(store.newArgumentVersionStatement.Close, errs)
	errs = mayAppendError(store.saveArgumentStatement.Close, errs)
	errs = mayAppendError(store.saveArgumentVersionStatement.Close, errs)
	errs = mayAppendError(store.saveClaimStatement.Close, errs)
	errs = mayAppendError(store.savePremiseStatement.Close, errs)
	if len(errs) == 0 {
		return nil
	}
	return closeErrors(errs)
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
