package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/wikisophia/api/server/arguments"
)

const deleteQuery = `UPDATE arguments SET deleted_on = $1 WHERE id = $2 RETURNING id;`

// Delete soft deletes an argument by ID.
func (store *PostgresStore) Delete(ctx context.Context, id int64) error {
	rows, err := store.deleteStatement.QueryContext(ctx, time.Now(), id)
	if err != nil {
		return err
	}
	defer tryClose(rows)

	hadRow := false
	for rows.Next() {
		hadRow = true
	}

	if hadRow {
		return nil
	}

	return &arguments.NotFoundError{
		Message: fmt.Sprintf("argument with id %d does not exist", id),
	}
}
