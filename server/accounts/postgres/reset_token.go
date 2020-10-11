package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/wikisophia/api/server/accounts"
	"github.com/wikisophia/api/server/accounts/tokens"
)

const newResetTokenQuery = `
WITH new_row AS (
	INSERT INTO accounts (email, reset_token, reset_token_expiry)
	SELECT $1, $2, $3
	WHERE NOT EXISTS (SELECT 1 FROM accounts WHERE email = $1)
	RETURNING id
)
SELECT id, true FROM new_row
UNION
SELECT id, false FROM accounts WHERE account = $1;
`

const updateResetTokenQuery = `
UPDATE accounts
SET reset_token = $2,
		reset_token_expiry = $3
		password_hash = NULL
WHERE email = $1
RETURNING id;
`

const resetTokenErrorMsg = "failed to save argument"

// See the docs on interfaces in store.go
func (store *PostgresStore) NewResetToken(ctx context.Context, email string) (accounts.Account, bool, error) {
	token, err := tokens.NewVerificationToken(50)
	expiration := time.Now().Add(24 * time.Hour)
	if err != nil {
		return accounts.Account{}, false, fmt.Errorf("%s: %v", resetTokenErrorMsg, err)
	}

	transaction, err := store.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return accounts.Account{}, false, fmt.Errorf("%s: %v", resetTokenErrorMsg, err)
	}
	row := transaction.QueryRow(ctx, newResetTokenQuery, email, token, expiration)
	var id int64
	var isNew bool
	if err := row.Scan(&id, &isNew); rollbackIfErr(ctx, transaction, err) {
		return accounts.Account{}, false, fmt.Errorf("%s: %v", resetTokenErrorMsg, err)
	}
	account := accounts.Account{
		ID:         id,
		Email:      email,
		ResetToken: token,
	}

	if !isNew {
		if result, err := transaction.Exec(ctx, updateResetTokenQuery, email, token, expiration); rollbackIfErr(ctx, transaction, err) {
			return accounts.Account{}, false, fmt.Errorf("%s: %v", resetTokenErrorMsg, err)
		} else if result.RowsAffected() != 1 {
			return account, false, nil
		}
	}

	if err := transaction.Commit(ctx); rollbackIfErr(ctx, transaction, err) {
		return accounts.Account{}, false, fmt.Errorf("%s: %v", resetTokenErrorMsg, err)
	}

	return account, true, nil
}

func rollbackIfErr(ctx context.Context, transaction pgx.Tx, err error) bool {
	if err != nil {
		if rollbackErr := transaction.Rollback(ctx); rollbackErr != nil {
			log.Printf("ERROR: Failed to rollback transaction: %v", rollbackErr)
		}
		return true
	}
	return false
}
