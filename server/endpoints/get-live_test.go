package endpoints_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wikisophia/api/server/arguments/argumentstest"
)

func TestGetLatest(t *testing.T) {
	expected := argumentstest.ParseSample(t, "../samples/save-request.json")
	var mistaken = expected
	mistaken.Premises = []string{"wrong", "stuff"}
	server := newAppForTests(t, nil).server
	id := doSaveObject(t, server, mistaken)
	expected.ID = id
	doValidUpdate(t, server, expected)
	expected.Version = 2
	rr := doGetArgument(server, id)
	assertSuccessfulJSON(t, rr)
	actual := parseArgumentResponse(t, rr.Body.Bytes())
	assert.Equal(t, expected, actual)
}

func TestGetMissingArgument(t *testing.T) {
	rr := doGetArgument(newAppForTests(t, nil).server, 1)
	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, "text/plain; charset=utf-8", rr.Header().Get("Content-Type"))
}

func TestGetStringID(t *testing.T) {
	rr := doRequest(newAppForTests(t, nil).server, httptest.NewRequest("GET", "/arguments/foo", nil))
	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, "text/plain; charset=utf-8", rr.Header().Get("Content-Type"))
}

func TestPostSpecificArgumentsNotAllowed(t *testing.T) {
	assertMethodNotAllowed(t, "POST", "/arguments/1")
	assertMethodNotAllowed(t, "POST", "/arguments/1/version/1")
}
