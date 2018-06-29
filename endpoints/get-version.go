package endpoints

import (
	"context"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Implements GET /arguments/:id/version/:version
func (s *Server) getArgumentVersion() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		id, goodID := parseIntParam(params.ByName("id"))
		versionInt, goodVersion := parseIntParam(params.ByName("version"))
		version, accurate := shrinkInt16(versionInt)

		// If the int can't fit into 16 bits, the database schema won't support it anyway.
		if !goodID || !goodVersion || !accurate {
			response := fmt.Sprintf("version %s of argument %s does not exist", params.ByName("version"), params.ByName("id"))
			http.Error(w, response, http.StatusNotFound)
			return
		}
		arg, err := s.argumentStore.FetchVersion(context.Background(), id, version)
		if writeStoreError(w, err) {
			return
		}
		writeArgument(w, arg, params.ByName("id"))
	}
}
