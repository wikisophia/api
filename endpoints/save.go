package endpoints

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/wikisophia/api-arguments/arguments"
)

// Implements POST /arguments
func (s *Server) saveArgument() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request: "+err.Error(), http.StatusBadRequest)
			return
		}

		var arg arguments.Argument
		if err := json.Unmarshal(data, &arg); err != nil {
			http.Error(w, "Failed to unmarshal argument: "+err.Error(), http.StatusBadRequest)
			return
		}
		if err := arguments.ValidateArgument(arg); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := s.argumentStore.Save(context.Background(), arg)
		if err != nil {
			http.Error(w, "Failed to save argument: "+err.Error(), http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Location", "/arguments/"+strconv.FormatInt(id, 10))

		// Firefox tries to parse the empty response to the AJAX call as XML
		// and throws an error without this
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
	}
}
