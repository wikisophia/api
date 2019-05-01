package endpoints

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Implements /suggestions
func suggestionsHandler() httprouter.Handle {
	var mockSuggestions = []byte(`["Socrates is mortal","Socrates is a man"]`)

	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		w.Write(mockSuggestions)
	}
}
