package accounts

import "context"

// Store combines all the functions needed to read & write Arguments
// into a single interface.
type Store interface {
	ResetTokenGenerator

	Close() error
}

// ResetTokenGenerator creates accounts and generate password reset tokens
type ResetTokenGenerator interface {
	// NewResetTokenWithAccount assigns a new password reset token to the account
	// with this email. If no accounts exist with this email, one will be created.
	NewResetTokenWithAccount(ctx context.Context, email string) (string, error)
}
