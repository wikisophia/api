package endpoints

import "github.com/julienschmidt/httprouter"

// newRouter sets up the service routes.
func newRouter(store Store) *httprouter.Router {
	router := httprouter.New()

	router.HandlerFunc("POST", "/arguments", saveHandler(store))
	router.HandlerFunc("GET", "/arguments", getAllArgumentsHandler(store))
	router.GET("/arguments/:id", getLiveArgumentHandler(store))
	router.PATCH("/arguments/:id", updateHandler(store))
	router.DELETE("/arguments/:id", deleteHandler(store))
	router.GET("/arguments/:id/version/:version", getArgumentByVersionHandler(store))

	return router
}

// setRoutes initializes the server with all its endpoints.
func (s *Server) setRoutes(store Store) {
	s.router.HandlerFunc("POST", "/arguments", saveHandler(store))
	s.router.HandlerFunc("GET", "/arguments", getAllArgumentsHandler(store))
	s.router.GET("/arguments/:id", getLiveArgumentHandler(store))
	s.router.PATCH("/arguments/:id", updateHandler(store))
	s.router.DELETE("/arguments/:id", deleteHandler(store))
	s.router.GET("/arguments/:id/version/:version", getArgumentByVersionHandler(store))
}
