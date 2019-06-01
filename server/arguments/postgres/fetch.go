package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"

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

// FetchSome returns all the "live" arguments matching the given options.
// If none exist, error will be nil and the slice empty.
func (store *Store) FetchSome(ctx context.Context, options arguments.FetchSomeOptions) ([]arguments.Argument, error) {
	selectArgumentsQuery := `SELECT arguments.id, arguments.live_version, argument_versions.id AS argument_version_id, claims.claim AS conclusion
	FROM arguments
		INNER JOIN argument_versions ON arguments.id = argument_versions.argument_id
		INNER JOIN claims ON argument_versions.conclusion_id = claims.id
	WHERE arguments.deleted = FALSE
		AND arguments.live_version = argument_versions.argument_version`

	var params []interface{}
	nextParamPlaceholder := newParamPlaceholderGenerator()
	if options.Conclusion != "" {
		selectArgumentsQuery += "\n\t\t AND claims.claim = " + nextParamPlaceholder()
		params = append(params, options.Conclusion)
	}
	selectArgumentsQuery += "\n\t ORDER BY arguments.id"
	if options.Count != 0 {
		selectArgumentsQuery += "\n\t LIMIT " + nextParamPlaceholder()
		params = append(params, options.Count)
	}
	if options.Offset != 0 {
		selectArgumentsQuery += "\n\t OFFSET " + nextParamPlaceholder()
		params = append(params, options.Offset)
	}

	fetchAllQuery := `WITH chosen_arguments AS (`
	fetchAllQuery += selectArgumentsQuery
	fetchAllQuery += ") \n"
	fetchAllQuery += `SELECT chosen_arguments.id, chosen_arguments.live_version, chosen_arguments.conclusion, claims.claim AS premise
	FROM chosen_arguments
		INNER JOIN argument_premises ON chosen_arguments.argument_version_id = argument_premises.argument_version_id
		INNER JOIN claims ON claims.id = argument_premises.premise_id;
	`

	rows, err := store.db.QueryContext(ctx, fetchAllQuery, params...)
	if err != nil {
		return nil, errors.Wrap(err, "failed fetchAll query")
	}
	defer tryClose(rows)

	args := make(map[int64]*arguments.Argument, 10)
	var id int64
	var version int
	var conclusion string
	var premise string
	for rows.Next() {
		if err := rows.Scan(&id, &version, &conclusion, &premise); err != nil {
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

// newParamPlaceholderGenerator returns a function that generates postgres wildcards.
// each time it's called, it returns a new one ($1, $2, $3, ...)
func newParamPlaceholderGenerator() func() string {
	index := 1
	return func() string {
		thisParam := "$" + strconv.Itoa(index)
		index++
		return thisParam
	}
}
