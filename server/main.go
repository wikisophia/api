package main

import (
	_ "net/http/pprof"

	"github.com/wikisophia/api-arguments/server/arguments/memory"
	argumentsInPostgres "github.com/wikisophia/api-arguments/server/arguments/postgres"
	"github.com/wikisophia/api-arguments/server/config"
	"github.com/wikisophia/api-arguments/server/endpoints"
	"github.com/wikisophia/api-arguments/server/postgres"
)

func main() {
	cfg := config.MustParse()
	server := endpoints.NewServer(*cfg.Server, newStore(cfg.Storage))

	done := make(chan struct{}, 1)
	go server.Start(done)
	<-done
}

func newStore(cfg *config.Storage) endpoints.Store {
	switch cfg.Type {
	case config.StorageTypeMemory:
		return memory.NewStore()
	case config.StorageTypePostgres:
		return argumentsInPostgres.NewStore(postgres.NewDB(cfg.Postgres))
	default:
		panic("Invalid config storage.type: " + cfg.Type + ". This should be caught during config valation.")
	}
}
