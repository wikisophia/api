package endpoints

import (
	"context"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Implements GET /arguments/:id
func (s *Server) getLiveArgument() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		id, goodID := parseIntParam(params.ByName("id"))
		if !goodID {
			http.Error(w, fmt.Sprintf("argument %s does not exist", params.ByName("id")), http.StatusNotFound)
			return
		}

		arg, err := s.argumentStore.FetchLive(context.Background(), int64(id))
		if writeStoreError(w, err) {
			return
		}
		writeArgument(w, arg, params.ByName("id"))
	}
}
