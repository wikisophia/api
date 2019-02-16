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
	SELECT argument_id, argument_version + 1, conclusion_id
		FROM argument_versions
		WHERE argument_id = $1
		ORDER BY argument_version DESC
		LIMIT 1
	RETURNING id, argument_version;
`

var updateLiveVersionQuery = `
UPDATE arguments
	SET live_version = $1
	WHERE id = $2;
`

const updateArgumentErrorMsg = "failed to update argument %d"

// UpdatePremises saves a new version of an argument with different premises.
func (store *Store) UpdatePremises(ctx context.Context, id int64, premises []string) (version int16, err error) {
	transaction, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return -1, errors.Wrap(err, updateArgumentErrorMsg)
	}
	argumentVersionID, argumentVersion, err := store.newArgumentVersion(ctx, transaction, id)
	if didRollback := rollbackIfErr(transaction, err); didRollback {
		return -1, err
	}
	err = store.savePremises(ctx, transaction, argumentVersionID, premises)
	if didRollback := rollbackIfErr(transaction, err); didRollback {
		return -1, errors.Wrapf(err, updateArgumentErrorMsg, id)
	}
	err = store.updateLiveVersion(ctx, transaction, argumentVersion, id)
	if didRollback := rollbackIfErr(transaction, err); didRollback {
		return -1, errors.Wrapf(err, updateArgumentErrorMsg, id)
	}
	err = transaction.Commit()
	if err != nil {
		return -1, errors.Wrapf(err, updateArgumentErrorMsg, id)
	}
	return argumentVersion, nil
}

func (store *Store) newArgumentVersion(ctx context.Context, transaction *sql.Tx, id int64) (int64, int16, error) {
	row := transaction.StmtContext(ctx, store.newArgumentVersionStatement).QueryRowContext(ctx, id)
	var argumentVersionID int64
	var argumentVersion int16
	if err := row.Scan(&argumentVersionID, &argumentVersion); err != nil {
		if err == sql.ErrNoRows {
			return -1, -1, &arguments.NotFoundError{
				Message: "argument " + strconv.FormatInt(id, 10) + " does not exist",
			}
		}
		return -1, -1, errors.Wrapf(err, `couldn't create new argument version for id=%d`, id)
	}
	return argumentVersionID, argumentVersion, nil
}

func (store *Store) updateLiveVersion(ctx context.Context, transaction *sql.Tx, liveVersion int16, argumentID int64) error {
	rows := transaction.StmtContext(ctx, store.updateLiveVersionStatement).QueryRowContext(ctx, liveVersion, argumentID)
	// Run Scan() to make sure the Rows get closed.
	err := rows.Scan()
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}
