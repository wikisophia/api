package postgres

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v4"
	"github.com/wikisophia/api/server/arguments"
)

var newArgumentVersionQuery = `
INSERT INTO argument_versions (argument_id, argument_version, conclusion_id)
	SELECT argument_id, argument_version + 1, $2
		FROM argument_versions
		WHERE argument_id = $1
		ORDER BY argument_version DESC
		LIMIT 1
	RETURNING id, argument_version;
`

const updateArgumentErrorMsg = "failed to update argument %d: %v"

// Update saves a new version of an argument.
func (store *PostgresStore) Update(ctx context.Context, argument arguments.Argument) (version int, err error) {
	tx, err := store.pool.BeginTx(ctx, pgx.TxOptions{})
	if didRollback := rollbackIfErr(ctx, tx, err); didRollback {
		return -1, err
	}
	conclusionID, err := store.saveClaim(ctx, tx, argument.Conclusion)
	if didRollback := rollbackIfErr(ctx, tx, err); didRollback {
		return -1, err
	}
	argumentVersionID, argumentVersion, err := store.newArgumentVersion(ctx, tx, argument.ID, conclusionID)
	if didRollback := rollbackIfErr(ctx, tx, err); didRollback {
		return -1, err
	}
	err = store.savePremises(ctx, tx, argumentVersionID, argument.Premises)
	if didRollback := rollbackIfErr(ctx, tx, err); didRollback {
		return -1, fmt.Errorf(updateArgumentErrorMsg, argument.ID, err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		return -1, fmt.Errorf(updateArgumentErrorMsg, argument.ID, err)
	}
	return argumentVersion, nil
}

func (store *PostgresStore) newArgumentVersion(ctx context.Context, tx pgx.Tx, argumentID int64, conclusionID int64) (int64, int, error) {
	row := tx.QueryRow(ctx, newArgumentVersionQuery, argumentID, conclusionID)
	var argumentVersionID int64
	var argumentVersion int
	if err := row.Scan(&argumentVersionID, &argumentVersion); err != nil {
		if err == pgx.ErrNoRows {
			return -1, -1, &arguments.NotFoundError{
				Message: "argument " + strconv.FormatInt(argumentID, 10) + " does not exist",
			}
		}
		return -1, -1, fmt.Errorf(`couldn't create new argument version for id=%d: %v`, argumentID, err)
	}
	return argumentVersionID, argumentVersion, nil
}
