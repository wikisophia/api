package main

import (
	_ "net/http/pprof"

	"github.com/hashicorp/go-multierror"
	"github.com/wikisophia/api/server/accounts"
	"github.com/wikisophia/api/server/arguments"
	"github.com/wikisophia/api/server/config"
	"github.com/wikisophia/api/server/endpoints"
	"github.com/wikisophia/api/server/postgres"
)

func main() {
	cfg := config.MustParse()
	store, cleanup := newStore(cfg.Storage)
	server := endpoints.NewServer(store)

	done := make(chan struct{}, 1)
	go server.Start(*cfg.Server, done)
	<-done
	cleanup()
}

func newStore(cfg *config.Storage) (endpoints.Store, func() error) {
	switch cfg.Type {
	case config.StorageTypeMemory:
		return memoryStore{
			accountsInMemory:  accounts.NewMemoryStore(),
			argumentsInMemory: arguments.NewMemoryStore(),
		}, func() error { return nil }
	case config.StorageTypePostgres:
		db := postgres.NewDB(cfg.Postgres)
		store := postgresStore{
			accountsInPostgres:  accounts.NewPostgresStore(db),
			argumentsInPostgres: arguments.NewPostgresStore(db),
		}
		return store, store.Close
	default:
		panic("Invalid config storage.type: " + cfg.Type + ". This should be caught during config valation.")
	}
}

type accountsInMemory = accounts.InMemoryStore
type argumentsInMemory = arguments.InMemoryStore
type memoryStore struct {
	*accountsInMemory
	*argumentsInMemory
}

type accountsInPostgres = accounts.PostgresStore
type argumentsInPostgres = arguments.PostgresStore
type postgresStore struct {
	*accountsInPostgres
	*argumentsInPostgres
}

func (store postgresStore) Close() error {
	var result *multierror.Error
	if err := store.accountsInPostgres.Close(); err != nil {
		multierror.Append(result, err)
	}
	if err := store.argumentsInPostgres.Close(); err != nil {
		multierror.Append(result, err)
	}
	return result.ErrorOrNil()
}
