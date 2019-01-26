package postgres

import "context"

const deleteQuery = `UPDATE arguments SET deleted = true WHERE id = $1;`

func (store *dbStore) Delete(ctx context.Context, id int64) error {
	rows, err := store.deleteStatement.QueryContext(ctx, id)
	if err != nil {
		return err
	}
	tryClose(rows)
	return nil
}
