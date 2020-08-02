package main

import (
	_ "net/http/pprof"

	"github.com/wikisophia/api-arguments/server/arguments"
	"github.com/wikisophia/api-arguments/server/config"
	"github.com/wikisophia/api-arguments/server/endpoints"
	"github.com/wikisophia/api-arguments/server/postgres"
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
		store := arguments.NewPostgresStore(postgres.NewDB(cfg.Postgres))
		return store, store.Close
	default:
		panic("Invalid config storage.type: " + cfg.Type + ". This should be caught during config valation.")
	}
}
