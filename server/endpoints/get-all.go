package endpoints

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

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

		exclude, ok := parseOptionalArrayOfInt64s(r.URL.Query().Get("exclude"))
		if !ok {
			http.Error(w, "The exclude query param must be a comma-separated list of non-negative integers.", http.StatusBadRequest)
			return
		}

		args, err := getter.FetchSome(context.Background(), arguments.FetchSomeOptions{
			Conclusion: r.URL.Query().Get("conclusion"),
			Count:      count,
			Exclude:    exclude,
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
	return parsed, true
}

func parseOptionalArrayOfInt64s(param string) ([]int64, bool) {
	if param == "" {
		return nil, true
	}
	stringNums := strings.Split(param, ",")
	intNums := make([]int64, 0, len(stringNums))
	for i := 0; i < len(stringNums); i++ {
		parsed, ok := parseInt64Param(stringNums[i])
		if !ok || parsed < 1 {
			return nil, false
		}
		intNums = append(intNums, parsed)
	}
	return intNums, true
}

func parseIntParam(param string) (int, bool) {
	parsed, err := strconv.Atoi(param)
	if err != nil || parsed < 1 {
		return 0, false
	}
	return parsed, true
}
