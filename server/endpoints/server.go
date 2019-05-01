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
		router: httprouter.New(),
	}

	// The HTTP Methods used here should stay in sync with the cors
	// AllowedMethods in Start().
	server.router.HandlerFunc("POST", "/arguments", saveHandler(store))
	server.router.HandlerFunc("GET", "/arguments", getAllArgumentsHandler(store))
	server.router.GET("/arguments/:id", getLiveArgumentHandler(store))
	server.router.PATCH("/arguments/:id", updateHandler(store))
	server.router.DELETE("/arguments/:id", deleteHandler(store))
	server.router.GET("/arguments/:id/version/:version", getArgumentByVersionHandler(store))
	server.router.GET("/suggestions", suggestionsHandler())
	return server
}

// Store has all the functions needed by the server for persistent storage
type Store interface {
	ArgumentDeleter
	ArgumentGetterByConclusion
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
	var handler http.Handler = s.router
	if len(s.config.CorsAllowedOrigins) > 0 {
		// AllowedMethods should stay in sync with the methods used by the routes
		handler = cors.New(cors.Options{
			AllowedOrigins: s.config.CorsAllowedOrigins,
			AllowedMethods: []string{"GET", "POST", "PATCH"},
			ExposedHeaders: []string{"Location"},
		}).Handler(handler)
	}

	httpServer := &http.Server{
		Addr:              s.config.Addr,
		Handler:           handler,
		ReadHeaderTimeout: s.config.ReadHeaderTimeout(),
	}

	go shutdownOnSignal(httpServer, done)
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
