package accounts

import "context"

// Store combines all the functions needed to read & write Arguments
// into a single interface.
type Store interface {
	Authenticator
	PasswordSetter
	ResetTokenGenerator
}
type Authenticator interface {
	// Authenticate returns the account's ID.
	//
	// If the email doesn't exist, it returns an AccountNotExistsError.
	// If the password is wrong, it returns an InvalidPasswordError.
	Authenticate(ctx context.Context, email, password string) (int64, error)
}

type PasswordSetter interface {
	// Change the password associated with the account using a reset token returned recently by
	//   ResetTokenGenerator.NewResetTokenWithAccount().
	//
	// If no account exists with this ID, it returns an AccountNotExistsError.
	// If the resetToken is wrong, it returns an InvalidResetTokenError.
	// If the password is unacceptable, it returns a ProhibitedPasswordError.
	SetForgottenPassword(ctx context.Context, id int64, password, resetToken string) error

	// Change the password for this account by using the old one, rather than a reset token.
	//
	// If the newPassword is unacceptable, it returns a ProhibitedPasswordError.
	// If no account with the ID exists, it returns an AccountNotExistsError.
	// If the old password is wrong, it returns an InvalidPasswordError.
	ChangePassword(ctx context.Context, id int64, oldPassword, newPassword string) error
}

type ResetTokenGenerator interface {
	// This associates a temporary password reset token with the account with the given email.
	// This token can be used in the PasswordSetter.SetForgottenPassword() method.
	//
	// If no Account exists with this email yet, one will be created. The bool return value is
	// true if the Account is new, and false if it existed already.
	NewResetToken(ctx context.Context, email string) (Account, bool, error)
}
