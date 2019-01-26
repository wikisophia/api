package endpoints

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"github.com/wikisophia/api-arguments/server/arguments"
	argumentsInMemory "github.com/wikisophia/api-arguments/server/arguments/memory"
	argumentsInPostgres "github.com/wikisophia/api-arguments/server/arguments/postgres"
	"github.com/wikisophia/api-arguments/server/config"
	"github.com/wikisophia/api-arguments/server/postgres"
)

// Server runs the service. Use NewServier() to construct one from an app config,
// and Start() to make it start listening and serving requests.
type Server struct {
	argumentStore arguments.Store
	config        *config.Configuration
	router        *httprouter.Router
}

// NewServer makes a server which defines REST endpoints for the service.
func NewServer(cfg config.Configuration) *Server {
	server := &Server{
		argumentStore: newArgumentStore(cfg.Storage),
		config:        &cfg,
		router:        httprouter.New(),
	}

	// The HTTP Methods used here should stay in sync with the cors
	// AllowedMethods in Start().
	server.router.HandlerFunc("POST", "/arguments", server.saveArgument())
	server.router.HandlerFunc("GET", "/arguments", server.getAllArguments())
	server.router.GET("/arguments/:id", server.getLiveArgument())
	server.router.PATCH("/arguments/:id", server.updateArgument())
	server.router.GET("/arguments/:id/version/:version", server.getArgumentVersion())

	// TODO: Suggestions should go into another service
	server.router.GET("/suggestions", server.suggestions())
	return server
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
	var handler http.Handler = s.router
	if len(s.config.Server.CorsAllowedOrigins) > 0 {
		// AllowedMethods should stay in sync with the methods used by the routes
		handler = cors.New(cors.Options{
			AllowedOrigins: s.config.Server.CorsAllowedOrigins,
			AllowedMethods: []string{"GET", "POST", "PATCH"},
			ExposedHeaders: []string{"Location"},
		}).Handler(handler)
	}

	httpServer := &http.Server{
		Addr:              s.config.Server.Addr,
		Handler:           handler,
		ReadHeaderTimeout: s.config.Server.ReadHeaderTimeout(),
	}

	go shutdownOnSignal(httpServer, done)
	return httpServer.ListenAndServe()
}

func newArgumentStore(cfg *config.Storage) arguments.Store {
	switch cfg.Type {
	case config.StorageTypeMemory:
		return argumentsInMemory.NewStore()
	case config.StorageTypePostgres:
		db := postgres.NewDB(cfg.Postgres)
		return argumentsInPostgres.NewStore(db)
	default:
		panic("Invalid config storage.type: " + cfg.Type + ". This should be caught during config valation.")
	}
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
