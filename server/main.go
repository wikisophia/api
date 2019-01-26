package main

import (
	_ "net/http/pprof"

	"github.com/wikisophia/api-arguments/server/config"
	"github.com/wikisophia/api-arguments/server/endpoints"
)

func main() {
	cfg := config.MustParse()

	done := make(chan struct{}, 1)
	server := endpoints.NewServer(cfg)
	go server.Start(done)

	<-done
}
