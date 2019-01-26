package postgres

import (
	"context"
	"database/sql"
	"log"

	"github.com/pkg/errors"
	"github.com/wikisophia/api-arguments/server/arguments"
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
INSERT INTO arguments (conclusion_id, latest_version, live_version)
	VALUES ($1, 1, 1)
	RETURNING id;
`

const savePremiseQuery = `
INSERT INTO premises (argument_id, argument_version, claim_id)
	VALUES ($1, $2, $3);
`

const saveArgumentErrorMsg = "failed to save argument"

func (store *dbStore) Save(ctx context.Context, argument arguments.Argument) (int64, error) {
	transaction, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return -1, errors.Wrap(err, saveArgumentErrorMsg)
	}
	conclusionID, err := store.saveClaim(ctx, transaction, argument.Conclusion)
	if didRollback := rollbackIfErr(transaction, err); didRollback {
		return -1, errors.Wrap(err, saveArgumentErrorMsg)
	}
	argumentID, err := store.saveArgument(ctx, transaction, conclusionID)
	if didRollback := rollbackIfErr(transaction, err); didRollback {
		return -1, errors.Wrap(err, saveArgumentErrorMsg)
	}
	err = store.savePremises(ctx, transaction, argumentID, 1, argument.Premises)
	if didRollback := rollbackIfErr(transaction, err); didRollback {
		return -1, errors.Wrap(err, saveArgumentErrorMsg)
	}
	err = transaction.Commit()
	if err != nil {
		return -1, errors.Wrap(err, saveArgumentErrorMsg)
	}
	return argumentID, nil
}

func (store *dbStore) saveClaim(ctx context.Context, tx *sql.Tx, claim string) (int64, error) {
	row := tx.StmtContext(ctx, store.saveClaimStatement).QueryRowContext(ctx, claim)
	var id int64
	if err := row.Scan(&id); err != nil {
		return -1, errors.Wrapf(err, `failed to save claim "%s"`, claim)
	}
	return id, nil
}

func (store *dbStore) saveArgument(ctx context.Context, tx *sql.Tx, conclusionID int64) (int64, error) {
	row := tx.StmtContext(ctx, store.saveArgumentStatement).QueryRowContext(ctx, conclusionID)
	var id int64
	if err := row.Scan(&id); err != nil {
		return -1, errors.Wrap(err, "failed to scan argument ID")
	}
	return id, nil
}

func (store *dbStore) savePremises(ctx context.Context, tx *sql.Tx, argumentID int64, argumentVersion int16, premises []string) error {
	for i := 0; i < len(premises); i++ {
		claimID, err := store.saveClaim(ctx, tx, premises[i])
		if err != nil {
			return errors.Wrapf(err, `failed to save premise as claim "%s"`, premises[i])
		}

		rows, err := tx.StmtContext(ctx, store.savePremiseStatement).QueryContext(ctx, argumentID, argumentVersion, claimID)
		if err != nil {
			return errors.Wrapf(err, `failed to save premise "%s"`, premises[i])
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
