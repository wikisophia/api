package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/pkg/errors"
	"github.com/wikisophia/api-arguments/arguments"
)

const fetchQuery = `
SELECT claims.claim, 'premise' as type
	FROM claims INNER JOIN premises ON claims.id = premises.claim_id
	WHERE premises.argument_id = $1 AND premises.argument_version = $2
UNION
SELECT claims.claim, 'conclusion' as type
	FROM claims INNER JOIN arguments ON claims.id = arguments.conclusion_id
	WHERE arguments.id = $1 AND arguments.deleted = false;
`

const fetchAllQuery = `
WITH these_arguments AS (
	SELECT arguments.id, arguments.live_version, arguments.isDefault
		FROM arguments INNER JOIN claims ON arguments.conclusion_id = claims.id
		WHERE claims.claim = $1
), these_premises AS (
	SELECT these_arguments.id, these_arguments.isDefault, premises.claim_id
		FROM these_arguments INNER JOIN premises ON these_arguments.live_version = premises.argument_version
		WHERE these_arguments.id = premises.argument_id
)
SELECT these_premises.id, these_premises.isDefault, claims.claim
	FROM these_premises INNER JOIN claims ON these_premises.claim_id = claims.id;
`

const fetchLiveVersionQuery = `SELECT live_version FROM arguments WHERE id = $1;`

func (store *dbStore) FetchVersion(ctx context.Context, id int64, version int16) (arguments.Argument, error) {
	rows, err := store.fetchStatement.QueryContext(ctx, id, version)
	if err != nil {
		return arguments.Argument{}, errors.Wrap(err, "argument fetch query failed")
	}
	defer tryClose(rows)

	var claim string
	var rowType string

	var conclusion string
	var premises []string

	for rows.Next() {
		if err := rows.Scan(&claim, &rowType); err != nil {
			return arguments.Argument{}, errors.Wrap(err, "fetch result scan failed")
		}
		switch rowType {
		case "premise":
			premises = append(premises, claim)
		case "conclusion":
			if conclusion != "" {
				return arguments.Argument{}, errors.Errorf("fetch returned two conclusions for version %d of argument %d, which suggests corrupt data", version, id)
			}
			conclusion = claim
		default:
			return arguments.Argument{}, errors.Errorf("fetch version %d for argument %d returned unknown error type: %s", version, id, rowType)
		}
	}
	if conclusion == "" {
		return arguments.Argument{}, &arguments.NotFoundError{
			Message: fmt.Sprintf("no argument found for version %d of argument %d", version, id),
		}
	}
	return arguments.Argument{
		Conclusion: conclusion,
		Premises:   premises,
	}, nil
}

func (store *dbStore) FetchLive(ctx context.Context, id int64) (arguments.Argument, error) {
	row := store.fetchLiveVersionStatement.QueryRowContext(ctx, id)
	var liveVersion int16
	if err := row.Scan(&liveVersion); err != nil {
		if err == sql.ErrNoRows {
			return arguments.Argument{}, &arguments.NotFoundError{
				Message: fmt.Sprintf("no argument found with id %d", id),
			}
		}

		return arguments.Argument{}, err
	}
	return store.FetchVersion(ctx, id, liveVersion)
}

func (store *dbStore) FetchAll(ctx context.Context, conclusion string) ([]arguments.ArgumentWithID, error) {
	rows, err := store.fetchAllStatement.QueryContext(ctx, conclusion)
	if err != nil {
		return nil, errors.Wrap(err, "failed fetchAll query")
	}
	defer tryClose(rows)

	args := make(map[int64]*arguments.ArgumentWithID, 10)
	var id int64
	var isDefault bool
	var premise string
	for rows.Next() {
		if err := rows.Scan(&id, &isDefault, &premise); err != nil {
			return nil, errors.Wrap(err, "fetch result scan failed")
		}
		if val, ok := args[id]; ok {
			val.Premises = append(val.Premises, premise)
		} else {
			premises := make([]string, 0, 10)
			premises = append(premises, premise)
			args[id] = &arguments.ArgumentWithID{
				Argument: arguments.Argument{
					Conclusion: conclusion,
					Premises:   premises,
				},
				ID: id,
			}
		}
	}
	toReturn := make([]arguments.ArgumentWithID, 0, len(args))
	for _, val := range args {
		toReturn = append(toReturn, *val)
	}
	return toReturn, nil
}

func tryClose(rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		log.Printf("ERROR: failed to close rows: %v", err)
	}
}
