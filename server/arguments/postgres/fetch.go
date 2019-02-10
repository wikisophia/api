package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/pkg/errors"
	"github.com/wikisophia/api-arguments/server/arguments"
)

const fetchQuery = `
SELECT claims.claim, 'premise' as type
FROM claims
	INNER JOIN argument_premises ON claims.id = argument_premises.premise_id
	INNER JOIN argument_versions ON argument_premises.argument_version_id = argument_versions.id
	INNER JOIN arguments ON arguments.id = argument_versions.argument_id
WHERE arguments.id = $1
	AND arguments.deleted = false
	AND argument_versions.argument_version = $2;
UNION
SELECT claims.claim, 'conclusion' as type
FROM claims
	INNER JOIN argument_versions ON claims.id = argument_versions.conclusion_id
	INNER JOIN arguments ON arguments.id = argument_versions.argument_id
WHERE arguments.id = $1
	AND arguments.deleted = false
	AND argument_versions.argument_version = $2;
`

const fetchAllQuery = `
SELECT arguments.id, premises.claim
	FROM arguments
		INNER JOIN argument_versions ON arguments.id = argument_versions.argument_id
		INNER JOIN argument_premises ON argument_versions.id = argument_premises.argument_version_id
		INNER JOIN claims conclusions ON argument_versions.conclusion_id = conclusions.id
		INNER JOIN claims premises ON argument_premises.premise_id = premises.id
	WHERE conclusions.claim = $1
		AND arguments.deleted = FALSE
		AND arguments.live_version = argument_versions.argument_version;
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

func (store *dbStore) FetchAll(ctx context.Context, conclusion string) ([]arguments.Argument, error) {
	rows, err := store.fetchAllStatement.QueryContext(ctx, conclusion)
	if err != nil {
		return nil, errors.Wrap(err, "failed fetchAll query")
	}
	defer tryClose(rows)

	args := make(map[int64]*arguments.Argument, 10)
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
			args[id] = &arguments.Argument{
				Conclusion: conclusion,
				Premises:   premises,
				ID:         id,
			}
		}
	}
	toReturn := make([]arguments.Argument, 0, len(args))
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
