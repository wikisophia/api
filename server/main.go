package main

import (
	_ "net/http/pprof"

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
		store := endpoints.AggregateStore{
			AccountsStore:  accounts.NewMemoryStore(),
			ArgumentsStore: arguments.NewMemoryStore(),
		}
		return store, store.Close
	case config.StorageTypePostgres:
		db := postgres.NewDB(cfg.Postgres)
		store := endpoints.AggregateStore{
			AccountsStore:  accounts.NewPostgresStore(db),
			ArgumentsStore: arguments.NewPostgresStore(db),
		}
		return store, store.Close
	default:
		panic("Invalid config storage.type: " + cfg.Type + ". This should be caught during config valation.")
	}
}
