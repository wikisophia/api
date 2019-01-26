package postgres

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/wikisophia/api-arguments/server/arguments"

	"github.com/pkg/errors"
)

var updateQuery = `
UPDATE arguments
	SET live_version = latest_version + 1, latest_version = latest_version + 1
	WHERE id = $1
	RETURNING latest_version;
`

const updateArgumentErrorMsg = "failed to update argument %d"

func (store *dbStore) UpdatePremises(ctx context.Context, id int64, premises []string) (version int16, err error) {
	transaction, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return -1, errors.Wrap(err, updateArgumentErrorMsg)
	}
	latestVersion, err := store.updateVersion(ctx, transaction, id)
	if didRollback := rollbackIfErr(transaction, err); didRollback {
		return -1, err
	}
	err = store.savePremises(ctx, transaction, id, latestVersion, premises)
	if didRollback := rollbackIfErr(transaction, err); didRollback {
		return -1, errors.Wrap(err, updateArgumentErrorMsg)
	}
	err = transaction.Commit()
	if err != nil {
		return -1, errors.Wrap(err, updateArgumentErrorMsg)
	}

	return latestVersion, nil
}

func (store *dbStore) updateVersion(ctx context.Context, transaction *sql.Tx, id int64) (int16, error) {
	row := transaction.StmtContext(ctx, store.updateStatement).QueryRowContext(ctx, id)
	var latestVersion int16
	if err := row.Scan(&latestVersion); err != nil {
		if err == sql.ErrNoRows {
			return -1, &arguments.NotFoundError{
				Message: "argument " + strconv.FormatInt(id, 10) + " does not exist",
			}
		}
		return -1, errors.Wrapf(err, `failed to update argument version %d`, id)
	}
	return latestVersion, nil
}
