package arguments

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

const saveClaimQuery = `
WITH new_row AS (
	INSERT INTO claims (claim)
	SELECT $1
	WHERE NOT EXISTS (SELECT 1 FROM claims WHERE claim = $1)
	RETURNING id
)
SELECT id FROM new_row
UNION
SELECT id FROM claims WHERE claim = $1;
`

const saveArgumentQuery = `
INSERT INTO arguments DEFAULT VALUES RETURNING id;
`

const saveArgumentVersionQuery = `
INSERT INTO argument_versions
	(argument_id, argument_version, conclusion_id) VALUES
	($1, 1, $2)
RETURNING id;
`

const savePremiseQuery = `
INSERT INTO argument_premises
	(argument_version_id, premise_id) VALUES
	($1, $2);
`

const saveArgumentErrorMsg = "failed to save argument"

// Save stores an argument and returns its ID.
// If the call succeeds, the Version will be 1.
func (store *Store) Save(ctx context.Context, argument Argument) (int64, error) {
	transaction, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return -1, fmt.Errorf("%s: %v", saveArgumentErrorMsg, err)
	}
	conclusionID, err := store.saveClaim(ctx, transaction, argument.Conclusion)
	if didRollback := rollbackIfErr(transaction, err); didRollback {
		return -1, fmt.Errorf("%s: %v", saveArgumentErrorMsg, err)
	}
	argumentID, err := store.saveArgument(ctx, transaction)
	if didRollback := rollbackIfErr(transaction, err); didRollback {
		return -1, fmt.Errorf("%s: %v", saveArgumentErrorMsg, err)
	}
	argumentVersionID, err := store.saveArgumentVersion(ctx, transaction, argumentID, 1, conclusionID)
	if didRollback := rollbackIfErr(transaction, err); didRollback {
		return -1, fmt.Errorf("%s: %v", saveArgumentErrorMsg, err)
	}

	err = store.savePremises(ctx, transaction, argumentVersionID, argument.Premises)
	if didRollback := rollbackIfErr(transaction, err); didRollback {
		return -1, fmt.Errorf("%s: %v", saveArgumentErrorMsg, err)
	}
	err = transaction.Commit()
	if err != nil {
		return -1, fmt.Errorf("%s: %v", saveArgumentErrorMsg, err)
	}
	return argumentID, nil
}

func (store *Store) saveClaim(ctx context.Context, tx *sql.Tx, claim string) (int64, error) {
	row := tx.StmtContext(ctx, store.saveClaimStatement).QueryRowContext(ctx, claim)
	var id int64
	if err := row.Scan(&id); err != nil {
		return -1, fmt.Errorf("failed to save claim \"%s\": %v", claim, err)
	}
	return id, nil
}

func (store *Store) saveArgument(ctx context.Context, tx *sql.Tx) (int64, error) {
	row := tx.StmtContext(ctx, store.saveArgumentStatement).QueryRowContext(ctx)
	var id int64
	if err := row.Scan(&id); err != nil {
		return -1, fmt.Errorf("failed to scan argument ID: %v", err)
	}
	return id, nil
}

func (store *Store) saveArgumentVersion(ctx context.Context, tx *sql.Tx, argumentID int64, versionID int, conclusionID int64) (int64, error) {
	row := tx.StmtContext(ctx, store.saveArgumentVersionStatement).QueryRowContext(ctx, argumentID, conclusionID)
	var id int64
	if err := row.Scan(&id); err != nil {
		return -1, fmt.Errorf("failed to scan argument ID: %v", err)
	}
	return id, nil
}

func (store *Store) savePremises(ctx context.Context, tx *sql.Tx, argumentVersionID int64, premises []string) error {
	for i := 0; i < len(premises); i++ {
		claimID, err := store.saveClaim(ctx, tx, premises[i])
		if err != nil {
			return fmt.Errorf(`failed to save premise as claim "%s": %v`, premises[i], err)
		}

		rows, err := tx.StmtContext(ctx, store.savePremiseStatement).QueryContext(ctx, argumentVersionID, claimID)
		if err != nil {
			return fmt.Errorf(`failed to save premise "%s": %v`, premises[i], err)
		}
		// Rows need to be closed before making the next query in a transaction.
		// See https://github.com/lib/pq/issues/81#issuecomment-229598201
		rows.Close()
	}
	return nil
}

func rollbackIfErr(transaction *sql.Tx, err error) bool {
	if err != nil {
		if rollbackErr := transaction.Rollback(); rollbackErr != nil {
			log.Printf("ERROR: Failed to rollback transaction: %v", rollbackErr)
		}
		return true
	}
	return false
}