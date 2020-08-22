package acceptancetest

import (
	"context"
	"errors"

	"github.com/wikisophia/api/server/accounts"
)

type Emailer struct {
	shouldSucceed bool

	Welcomes       []*accounts.Account
	PasswordResets []*accounts.Account
}

func (e *Emailer) SendWelcome(ctx context.Context, account accounts.Account) error {
	e.Welcomes = append(e.Welcomes, &account)
	if e.shouldSucceed {
		return nil
	}
	return errors.New("Welcome message failed to send")
}

func (e *Emailer) SendReset(ctx context.Context, account accounts.Account) error {
	e.PasswordResets = append(e.PasswordResets, &account)
	if e.shouldSucceed {
		return nil
	}
	return errors.New("Password reset message failed to send")
}
