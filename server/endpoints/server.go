package endpoints

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/julienschmidt/httprouter"
	"github.com/wikisophia/api-arguments/server/config"
)

// Server runs the service. Use NewServier() to construct one from an app config,
// and Start() to make it start listening and serving requests.
type Server struct {
	config *config.Server
	router *httprouter.Router
}

// NewServer makes a server which defines REST endpoints for the service.
func NewServer(cfg config.Server, store Store) *Server {
	server := &Server{
		config: &cfg,
		router: newRouter(store),
	}
	return server
}

// Store has all the functions needed by the server for persistent storage
type Store interface {
	ArgumentDeleter
	ArgumentsGetter
	ArgumentGetterByVersion
	ArgumentGetterLiveVersion
	ArgumentSaver
	ArgumentUpdater
}

// Handle exists to make testing easier.
// It lets the Server act without having to bind to a port.
func (s *Server) Handle(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}

// Start connects the API server to its port and blocks until it hears a
// shutdown signal. Once the server has shut down completely, it adds
// an element to the done channel.
func (s *Server) Start(done chan<- struct{}) error {
	httpServer := &http.Server{
		Addr:              s.config.Addr,
		Handler:           s.router,
		ReadHeaderTimeout: s.config.ReadHeaderTimeout(),
	}

	go shutdownOnSignal(httpServer, done)
	if s.config.UseSSL {
		return httpServer.ListenAndServeTLS(s.config.CertPath, s.config.KeyPath)
	}
	return httpServer.ListenAndServe()
}

func shutdownOnSignal(server *http.Server, done chan<- struct{}) {
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)
	sig := <-signals
	log.Printf("Received signal %v. API server shutting down.", sig)
	server.Shutdown(context.Background())
	var s struct{}
	done <- s
}
