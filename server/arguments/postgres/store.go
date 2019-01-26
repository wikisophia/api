package postgres

import (
	"database/sql"
	"log"
	"strings"

	"github.com/wikisophia/api-arguments/server/arguments"
)

// NewStore returns a Store which is used to save and load Arguments.
// The db should point to a Postgres database.
// The returned Store.Close() function will *not* close this connection, since we did not open it.
func NewStore(db *sql.DB) arguments.Store {
	if db == nil {
		log.Fatalf("A database connection is required to make a Store.")
	}
	return &dbStore{
		db:                        db,
		deleteStatement:           mustPrepareQuery(db, deleteQuery),
		fetchAllStatement:         mustPrepareQuery(db, fetchAllQuery),
		fetchStatement:            mustPrepareQuery(db, fetchQuery),
		fetchLiveVersionStatement: mustPrepareQuery(db, fetchLiveVersionQuery),
		saveArgumentStatement:     mustPrepareQuery(db, saveArgumentQuery),
		saveClaimStatement:        mustPrepareQuery(db, saveClaimQuery),
		savePremiseStatement:      mustPrepareQuery(db, savePremiseQuery),
		updateStatement:           mustPrepareQuery(db, updateQuery),
	}
}

// The dbStore expects that {projectRoot}/postgres/scripts/init.sql
// has already been run on your database so that the expected schema exists.
type dbStore struct {
	db                        *sql.DB
	deleteStatement           *sql.Stmt
	fetchAllStatement         *sql.Stmt
	fetchStatement            *sql.Stmt
	fetchLiveVersionStatement *sql.Stmt
	saveArgumentStatement     *sql.Stmt
	saveClaimStatement        *sql.Stmt
	savePremiseStatement      *sql.Stmt
	updateStatement           *sql.Stmt
}

type closeErrors []error

func (errs closeErrors) Error() string {
	if len(errs) == 0 {
		return ""
	}

	sb := strings.Builder{}
	sb.WriteString("error(s) occurred while closing the Store:\n")
	for i := 0; i < len(errs); i++ {
		sb.WriteString("  ")
		sb.WriteString(errs[i].Error())
		sb.WriteString("\n")
	}
	return sb.String()
}

func (store *dbStore) Close() error {
	var errs []error
	errs = mayAppendError(store.deleteStatement.Close, errs)
	errs = mayAppendError(store.fetchAllStatement.Close, errs)
	errs = mayAppendError(store.fetchStatement.Close, errs)
	errs = mayAppendError(store.fetchLiveVersionStatement.Close, errs)
	errs = mayAppendError(store.saveArgumentStatement.Close, errs)
	errs = mayAppendError(store.saveClaimStatement.Close, errs)
	errs = mayAppendError(store.savePremiseStatement.Close, errs)
	errs = mayAppendError(store.updateStatement.Close, errs)
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
		log.Fatalf("Failed to prepare statement with query %s. Error was %v", saveArgumentQuery, err)
	}
	return statement
}
