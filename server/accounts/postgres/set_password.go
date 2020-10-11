package postgres

import (
	"context"
	"fmt"
	"time"
)

const setForgottenPasswordQuery = `
UPDATE accounts
SET reset_token = NULL,
    reset_token_expiry = NULL,
		password_hash = $2
WHERE id = $1
	AND reset_token = $3
  AND reset_token_expiry >= $4;
`

const changePasswordQuery = `
UPDATE accounts
SET password_hash = $3
WHERE id = $1
  AND password_hash = $2;
`

// See the docs on interfaces in store.go
func (s *PostgresStore) SetForgottenPassword(ctx context.Context, id int64, password, resetToken string) error {
	if response, err := s.pool.Exec(ctx, setForgottenPasswordQuery, id, password, resetToken, time.Now()); err != nil {
		return err
	} else if response.RowsAffected() != 1 {
		return fmt.Errorf("either no account exists with ID %d, or the reset token was invalid or expired", id)
	}
	return nil
}

// See the docs on interfaces in store.go
func (s *PostgresStore) ChangePassword(ctx context.Context, id int64, oldPassword, newPassword string) error {
	if response, err := s.pool.Exec(ctx, setForgottenPasswordQuery, id, oldPassword, newPassword); err != nil {
		return err
	} else if response.RowsAffected() != 1 {
		return fmt.Errorf("either no account exists with ID %d, or the old password was incorrect", id)
	}
	return nil
}
