package postgres

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/wikisophia/api-arguments/server/arguments"

	"github.com/pkg/errors"
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

const updateArgumentErrorMsg = "failed to update argument %d"

// Update saves a new version of an argument.
func (store *Store) Update(ctx context.Context, argument arguments.Argument) (version int, err error) {
	transaction, err := store.db.BeginTx(ctx, nil)
	if didRollback := rollbackIfErr(transaction, err); didRollback {
		return -1, err
	}
	conclusionID, err := store.saveClaim(ctx, transaction, argument.Conclusion)
	if didRollback := rollbackIfErr(transaction, err); didRollback {
		return -1, err
	}
	argumentVersionID, argumentVersion, err := store.newArgumentVersion(ctx, transaction, argument.ID, conclusionID)
	if didRollback := rollbackIfErr(transaction, err); didRollback {
		return -1, err
	}
	err = store.savePremises(ctx, transaction, argumentVersionID, argument.Premises)
	if didRollback := rollbackIfErr(transaction, err); didRollback {
		return -1, errors.Wrapf(err, updateArgumentErrorMsg, argument.ID)
	}
	err = transaction.Commit()
	if err != nil {
		return -1, errors.Wrapf(err, updateArgumentErrorMsg, argument.ID)
	}
	return argumentVersion, nil
}

func (store *Store) newArgumentVersion(ctx context.Context, transaction *sql.Tx, argumentID int64, conclusionID int64) (int64, int, error) {
	row := transaction.StmtContext(ctx, store.newArgumentVersionStatement).QueryRowContext(ctx, argumentID, conclusionID)
	var argumentVersionID int64
	var argumentVersion int
	if err := row.Scan(&argumentVersionID, &argumentVersion); err != nil {
		if err == sql.ErrNoRows {
			return -1, -1, &arguments.NotFoundError{
				Message: "argument " + strconv.FormatInt(argumentID, 10) + " does not exist",
			}
		}
		return -1, -1, errors.Wrapf(err, `couldn't create new argument version for id=%d`, argumentID)
	}
	return argumentVersionID, argumentVersion, nil
}
