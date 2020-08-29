package main

import (
	_ "net/http/pprof"

	"github.com/wikisophia/api/server/accounts"
	"github.com/wikisophia/api/server/accounts/email"
	accountsMemory "github.com/wikisophia/api/server/accounts/memory"
	accountsPostgres "github.com/wikisophia/api/server/accounts/postgres"
	"github.com/wikisophia/api/server/arguments"
	argumentsMemory "github.com/wikisophia/api/server/arguments/memory"
	argumentsPostgres "github.com/wikisophia/api/server/arguments/postgres"
	"github.com/wikisophia/api/server/http"

	"github.com/wikisophia/api/server/config"
	"github.com/wikisophia/api/server/postgres"
)

func main() {
	cfg := config.MustParse()
	deps := newDependencies(&cfg)
	server := http.NewServer(cfg.JwtPrivateKey(), deps)

	done := make(chan struct{}, 1)
	go server.Start(*cfg.Server, done)
	<-done
}

func newDependencies(cfg *config.Configuration) http.Dependencies {
	return http.ServerDependencies{
		AccountsStore:  newAccountsStore(cfg.AccountsStore),
		ArgumentsStore: newArgumentsStore(cfg.ArgumentsStore),
		Emailer:        email.ConsoleEmailer{},
	}
}

func newAccountsStore(cfg *config.Storage) accounts.Store {
	switch cfg.Type {
	case config.StorageTypeMemory:
		return accountsMemory.NewMemoryStore()
	case config.StorageTypePostgres:
		return accountsPostgres.NewPostgresStore(postgres.NewPGXPool(cfg.Postgres))
	default:
		panic("Invalid config storage.type: " + cfg.Type + ". This should be caught during config valation.")
	}
}

func newArgumentsStore(cfg *config.Storage) arguments.Store {
	switch cfg.Type {
	case config.StorageTypeMemory:
		return argumentsMemory.NewMemoryStore()
	case config.StorageTypePostgres:
		return argumentsPostgres.NewPostgresStore(postgres.NewPGXPool(cfg.Postgres))
	default:
		panic("Invalid config storage.type: " + cfg.Type + ". This should be caught during config valation.")
	}
}
