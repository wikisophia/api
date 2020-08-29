package main

import (
	_ "net/http/pprof"

	"github.com/wikisophia/api/server/accounts/email"
	accountsMemory "github.com/wikisophia/api/server/accounts/memory"
	accountsPostgres "github.com/wikisophia/api/server/accounts/postgres"
	argumentsMemory "github.com/wikisophia/api/server/arguments/memory"
	argumentsPostgres "github.com/wikisophia/api/server/arguments/postgres"
	"github.com/wikisophia/api/server/http"

	"github.com/wikisophia/api/server/config"
	"github.com/wikisophia/api/server/postgres"
)

func main() {
	cfg := config.MustParse()
	store, cleanup := newDependencies(cfg.Storage)
	server := http.NewServer(cfg.JwtPrivateKey(), store)

	done := make(chan struct{}, 1)
	go server.Start(*cfg.Server, done)
	<-done
	cleanup()
}

func newDependencies(cfg *config.Storage) (http.Dependencies, func() error) {
	switch cfg.Type {
	case config.StorageTypeMemory:
		store := http.ServerDependencies{
			AccountsStore:  accountsMemory.NewMemoryStore(),
			ArgumentsStore: argumentsMemory.NewMemoryStore(),
			Emailer:        email.ConsoleEmailer{},
		}
		return store, store.Close
	case config.StorageTypePostgres:
		pool := postgres.NewPGXPool(cfg.Postgres)
		store := http.ServerDependencies{
			AccountsStore:  accountsPostgres.NewPostgresStore(pool),
			ArgumentsStore: argumentsPostgres.NewPostgresStore(pool),
			Emailer:        email.ConsoleEmailer{},
		}
		return store, store.Close
	default:
		panic("Invalid config storage.type: " + cfg.Type + ". This should be caught during config valation.")
	}
}
