package email

import (
	"context"
	"log"

	"github.com/wikisophia/api/server/accounts"
)

type ConsoleEmailer struct{}

func (e ConsoleEmailer) SendWelcome(ctx context.Context, account accounts.Account) error {
	log.Printf("%s has ID %d and reset token %s", account.Email, account.ID, account.ResetToken)
	return nil
}
func (e ConsoleEmailer) SendReset(ctx context.Context, account accounts.Account) error {
	log.Printf("%s has ID %d and new reset token %s", account.Email, account.ID, account.ResetToken)
	return nil
}
