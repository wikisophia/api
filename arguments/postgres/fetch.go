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

func tryClose(rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		log.Printf("ERROR: failed to close rows: %v", err)
	}
}
