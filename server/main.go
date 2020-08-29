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
	store := newDependencies(cfg.Storage)
	server := http.NewServer(cfg.JwtPrivateKey(), store)

	done := make(chan struct{}, 1)
	go server.Start(*cfg.Server, done)
	<-done
}

func newDependencies(cfg *config.Storage) http.Dependencies {
	switch cfg.Type {
	case config.StorageTypeMemory:
		return http.ServerDependencies{
			AccountsStore:  accountsMemory.NewMemoryStore(),
			ArgumentsStore: argumentsMemory.NewMemoryStore(),
			Emailer:        email.ConsoleEmailer{},
		}
	case config.StorageTypePostgres:
		pool := postgres.NewPGXPool(cfg.Postgres)
		return http.ServerDependencies{
			AccountsStore:  accountsPostgres.NewPostgresStore(pool),
			ArgumentsStore: argumentsPostgres.NewPostgresStore(pool),
			Emailer:        email.ConsoleEmailer{},
		}
	default:
		panic("Invalid config storage.type: " + cfg.Type + ". This should be caught during config valation.")
	}
}
