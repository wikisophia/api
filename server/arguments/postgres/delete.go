package postgres

import (
	"context"
	"fmt"

	"github.com/wikisophia/api-arguments/server/arguments"
)

const deleteQuery = `UPDATE arguments SET deleted = true WHERE id = $1 RETURNING id;`

func (store *dbStore) Delete(ctx context.Context, id int64) error {
	rows, err := store.deleteStatement.QueryContext(ctx, id)
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
