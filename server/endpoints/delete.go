package endpoints

import (
	"context"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s *Server) deleteArgument() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		id, goodID := parseIntParam(params.ByName("id"))
		if !goodID {
			http.Error(w, fmt.Sprintf("argument %s does not exist", params.ByName("id")), http.StatusNotFound)
			return
		}
		if err := s.argumentStore.Delete(context.Background(), id); writeStoreError(w, err) {
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNoContent)
	}
}
