package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
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

const selectPasswordByIdQuery = `
SELECT password_hash
FROM accounts
WHERE id = $1;
`

const changePasswordQuery = `
UPDATE accounts
SET password_hash = $2
WHERE id = $1;
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
	row := s.pool.QueryRow(ctx, selectPasswordByIdQuery, id)
	var oldPasswordHash string
	if err := row.Scan(&oldPasswordHash); err == pgx.ErrNoRows {
		return fmt.Errorf("either no account exists with ID %d, or the old password was incorrect", id)
	} else if err != nil {
		return fmt.Errorf("failed to change password: %v", err)
	}
	matches, err := s.hasher.Matches(oldPassword, oldPasswordHash)
	if err != nil {
		return errors.New("error matching password against the database")
	}
	if !matches {
		return fmt.Errorf("either no account exists with ID %d, or the old password was incorrect", id)
	}
	newHash, err := s.hasher.Hash(newPassword)
	if err != nil {
		return fmt.Errorf("failed to change password: %v", err)
	}
	_, err = s.pool.Exec(ctx, changePasswordQuery, id, newHash)
	if err != nil {
		return fmt.Errorf("failed to change password: %v", err)
	}
	return nil
}
