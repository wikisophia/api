package accounts

import "context"

// Emailer sends new account and password reset tokens
type Emailer interface {
	SendWelcome(ctx context.Context, account Account) error
	SendReset(ctx context.Context, account Account) error
}
