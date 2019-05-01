package endpoints

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/wikisophia/api-arguments/server/arguments"
)

// ArgumentGetterByConclusion can fetch all the arguments with a given conclusion.
type ArgumentGetterByConclusion interface {
	// FetchAll finds all the available arguments for a conclusion.
	// If none exist, error will be nil and the slice empty.
	FetchAll(ctx context.Context, conclusion string) ([]arguments.Argument, error)
}

func getAllArgumentsHandler(getter ArgumentGetterByConclusion) http.HandlerFunc {
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

		args, err := getter.FetchAll(context.Background(), conclusion)
		if err != nil {
			http.Error(w, "failed to fetch arguments from the backend", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		log.Printf("returning: %v", GetAllResponse{
			Arguments: args,
		})
		if err = json.NewEncoder(w).Encode(GetAllResponse{
			Arguments: args,
		}); err != nil {
			log.Printf("ERROR: Failed encoding response to GET /arguments for conclusion \"%s\": %v", conclusion, err)
		}
	}
}

// GetAllResponse is the contract class for the GET /arguments?conclusion=foo endpoint
type GetAllResponse struct {
	Arguments []arguments.Argument `json:"arguments"`
}
