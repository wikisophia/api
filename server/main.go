package main

import (
	_ "net/http/pprof"

	"github.com/wikisophia/api-arguments/arguments"
	"github.com/wikisophia/api-arguments/config"
	"github.com/wikisophia/api-arguments/endpoints"
	"github.com/wikisophia/api-arguments/postgres"
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
		return arguments.NewMemoryStore(), func() error { return nil }
	case config.StorageTypePostgres:
		store := postgres.NewStore(postgres.NewDB(cfg.Postgres))
		return store, store.Close
	default:
		panic("Invalid config storage.type: " + cfg.Type + ". This should be caught during config valation.")
	}
}
