package endpoints

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/wikisophia/api-arguments/server/arguments"
)

// ArgumentsGetter can fetch lists of arguments at once.
type ArgumentsGetter interface {
	// FetchSome finds the arguments which match the options.
	// If none exist, error will be nil and the slice empty.
	FetchSome(ctx context.Context, options arguments.FetchSomeOptions) ([]arguments.Argument, error)
}

func getAllArgumentsHandler(getter ArgumentsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL == nil {
			http.Error(w, "URL was nil. Bad Request-Line?", http.StatusBadRequest)
			return
		}
		count, ok := parseOptionalNonNegativeIntParam(r.URL.Query().Get("count"))
		if !ok {
			http.Error(w, "The count query param must be a nonnegative integer.", http.StatusBadRequest)
			return
		}
		offset, ok := parseOptionalNonNegativeIntParam(r.URL.Query().Get("offset"))
		if !ok {
			http.Error(w, "The offset query param must be a nonnegative integer.", http.StatusBadRequest)
			return
		}

		args, err := getter.FetchSome(context.Background(), arguments.FetchSomeOptions{
			Conclusion: r.URL.Query().Get("conclusion"),
			Count:      int(count),
			Offset:     offset,
		})
		if err != nil {
			http.Error(w, "failed to fetch arguments from the backend", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(GetAllResponse{
			Arguments: args,
		})
	}
}

// GetAllResponse is the contract class for the GET /arguments?conclusion=foo endpoint
type GetAllResponse struct {
	Arguments []arguments.Argument `json:"arguments"`
}

func parseOptionalNonNegativeIntParam(param string) (int, bool) {
	if param == "" {
		return 0, true
	}
	parsed, ok := parseIntParam(param)
	if !ok || parsed < 1 {
		return 0, false
	}
	shrunk, ok := shrinkInt(parsed)
	if !ok {
		return 0, false
	}
	return shrunk, true
}
