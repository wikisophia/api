package http_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wikisophia/api/server/acceptancetest"
)

func TestGetLatest(t *testing.T) {
	expected := acceptancetest.ParseSample(t, samplesPath+"save-request.json")
	var mistaken = expected
	mistaken.Premises = []string{"wrong", "stuff"}
	app := newApp(t, nil)
	id := app.SaveSuccessfully(t, mistaken)
	expected.ID = id
	app.UpdateSuccessfully(t, expected)
	expected.Version = 2
	actual := app.GetLiveSuccessfully(id)
	assert.Equal(t, expected, actual)
}

func TestGetMissingArgument(t *testing.T) {
	rr := newApp(t, nil).Do(newGetArgument(1))
	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, "text/plain; charset=utf-8", rr.Header().Get("Content-Type"))
}

func TestGetStringID(t *testing.T) {
	rr := newApp(t, nil).Do(httptest.NewRequest("GET", "/arguments/foo", nil))
	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, "text/plain; charset=utf-8", rr.Header().Get("Content-Type"))
}

func TestPostSpecificArgumentsNotAllowed(t *testing.T) {
	acceptancetest.AssertMethodNotAllowed(t, "POST", "/arguments/1")
	acceptancetest.AssertMethodNotAllowed(t, "POST", "/arguments/1/version/1")
}

func newGetArgument(id int64) *http.Request {
	return httptest.NewRequest("GET", "/arguments/"+strconv.FormatInt(id, 10), nil)
}
