package main

import (
	_ "net/http/pprof"

	"github.com/wikisophia/api-arguments/config"
	"github.com/wikisophia/api-arguments/endpoints"
)

func main() {
	cfg := config.MustParseConfig()

	done := make(chan struct{}, 1)
	server := endpoints.NewServer(cfg)
	go server.Start(done)

	<-done
}
