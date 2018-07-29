package endpoints

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

func (s *Server) getAllArguments() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL == nil {
			http.Error(w, "URL was nil. Bad Request-Line?", http.StatusBadRequest)
			return
		}
		conclusion := r.URL.Query().Get("conclusion")
		if conclusion == "" {
			http.Error(w, "missing required query parameter: conclusion", http.StatusBadRequest)
			return
		}

		args, err := s.argumentStore.FetchAll(context.Background(), conclusion)
		if err != nil {
			http.Error(w, "failed to fetch arguments from the backend", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(args); err != nil {
			log.Printf("ERROR: Failed encoding response to GET /arguments for conclusion \"%s\": %v", conclusion, err)
		}
	}
}
