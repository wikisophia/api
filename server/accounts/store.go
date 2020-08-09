package accounts

import "context"

// Store combines all the functions needed to read & write Arguments
// into a single interface.
type Store interface {
	Authenticator
	PasswordSetter
	ResetTokenGenerator
	Close() error
}

type Authenticator interface {
	// Authenticate returns the account's ID.
	// If the email doesn't exist, it returns an EmailNotExistsError.
	// If the password is wrong, it returns an InvalidPasswordError.
	Authenticate(ctx context.Context, email, password string) (int64, error)
}

type PasswordSetter interface {
	// SetPassword changes the password associated with the email and returns the account's ID.
	// If the email doesn't exist, it returns an EmailNotExistsError.
	// If the resetToken is wrong (expired or never returned by ResetPassword(email)),
	//   it returns an InvalidPasswordError.
	SetPassword(ctx context.Context, email, password, resetToken string) (int64, error)
}

type ResetTokenGenerator interface {
	// NewResetTokenWithAccount assigns a new password reset token to the account
	// with this email. If no accounts exist with this email, one will be created.
	NewResetTokenWithAccount(ctx context.Context, email string) (string, error)
}
