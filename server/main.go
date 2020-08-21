package main

import (
	"context"
	"log"
	_ "net/http/pprof"

	"github.com/wikisophia/api/server/accounts"
	"github.com/wikisophia/api/server/arguments"
	"github.com/wikisophia/api/server/config"
	"github.com/wikisophia/api/server/endpoints"
	"github.com/wikisophia/api/server/postgres"
)

func main() {
	cfg := config.MustParse()
	store, cleanup := newDependencies(cfg.Storage)
	server := endpoints.NewServer(cfg.JwtPrivateKey(), store)

	done := make(chan struct{}, 1)
	go server.Start(*cfg.Server, done)
	<-done
	cleanup()
}

func newDependencies(cfg *config.Storage) (endpoints.Dependencies, func() error) {
	switch cfg.Type {
	case config.StorageTypeMemory:
		store := endpoints.ServerDependencies{
			AccountsStore:  accounts.NewMemoryStore(),
			ArgumentsStore: arguments.NewMemoryStore(),
			Emailer:        ConsoleEmailer{},
		}
		return store, store.Close
	case config.StorageTypePostgres:
		db := postgres.NewDB(cfg.Postgres)
		store := endpoints.ServerDependencies{
			AccountsStore:  accounts.NewPostgresStore(db),
			ArgumentsStore: arguments.NewPostgresStore(db),
			Emailer:        ConsoleEmailer{},
		}
		return store, store.Close
	default:
		panic("Invalid config storage.type: " + cfg.Type + ". This should be caught during config valation.")
	}
}

type ConsoleEmailer struct{}

func (e ConsoleEmailer) SendWelcome(ctx context.Context, account accounts.Account) error {
	log.Printf("%s has ID %d and reset token %s", account.Email, account.ID, account.ResetToken)
	return nil
}
func (e ConsoleEmailer) SendReset(ctx context.Context, account accounts.Account) error {
	log.Printf("%s has ID %d and new reset token %s", account.Email, account.ID, account.ResetToken)
	return nil
}
