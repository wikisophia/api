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
(SELECT claims.claim, argument_premises.id AS o
FROM claims
	INNER JOIN argument_premises ON claims.id = argument_premises.premise_id
	INNER JOIN argument_versions ON argument_premises.argument_version_id = argument_versions.id
	INNER JOIN arguments ON arguments.id = argument_versions.argument_id
WHERE arguments.id = $1
	AND arguments.deleted = false
	AND argument_versions.argument_version = $2)
UNION ALL
(SELECT claims.claim, -1 AS o
FROM claims
	INNER JOIN argument_versions ON claims.id = argument_versions.conclusion_id
	INNER JOIN arguments ON arguments.id = argument_versions.argument_id
WHERE arguments.id = $1
	AND arguments.deleted = false
	AND argument_versions.argument_version = $2)
ORDER BY o;
`

const fetchAllQuery = `
SELECT arguments.id, arguments.live_version, premises.claim
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

// FetchVersion fetches a specific version of an argument.
func (store *Store) FetchVersion(ctx context.Context, id int64, version int) (arguments.Argument, error) {
	rows, err := store.fetchStatement.QueryContext(ctx, id, version)
	if err != nil {
		return arguments.Argument{}, errors.Wrap(err, "argument fetch query failed")
	}
	defer tryClose(rows)

	var claim string
	var dummy int64

	var conclusion string
	var premises []string

	for rows.Next() {
		if err := rows.Scan(&claim, &dummy); err != nil {
			return arguments.Argument{}, errors.Wrap(err, "fetch result scan failed")
		}
		if conclusion == "" {
			conclusion = claim
		} else {
			premises = append(premises, claim)
		}
	}
	if conclusion == "" {
		return arguments.Argument{}, &arguments.NotFoundError{
			Message: fmt.Sprintf("no argument found for version %d of argument %d", version, id),
		}
	}
	return arguments.Argument{
		ID:         id,
		Version:    version,
		Conclusion: conclusion,
		Premises:   premises,
	}, nil
}

// FetchLive fetches the "active" version of an argument.
// This is usually the newest one, but it may not be if an
// update has been reverted.
func (store *Store) FetchLive(ctx context.Context, id int64) (arguments.Argument, error) {
	row := store.fetchLiveVersionStatement.QueryRowContext(ctx, id)
	var liveVersion int
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

// FetchAll returns all the "live" arguments with a given conclusion.
func (store *Store) FetchAll(ctx context.Context, conclusion string) ([]arguments.Argument, error) {
	rows, err := store.fetchAllStatement.QueryContext(ctx, conclusion)
	if err != nil {
		return nil, errors.Wrap(err, "failed fetchAll query")
	}
	defer tryClose(rows)

	args := make(map[int64]*arguments.Argument, 10)
	var id int64
	var version int
	var premise string
	for rows.Next() {
		if err := rows.Scan(&id, &version, &premise); err != nil {
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
				Version:    version,
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
