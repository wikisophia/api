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
