package email

import (
	"context"

	"github.com/wikisophia/api/server/accounts"
)

// Emailer sends new account and password reset tokens
type Emailer interface {
	SendWelcome(ctx context.Context, account accounts.Account) error
	SendReset(ctx context.Context, account accounts.Account) error
}
