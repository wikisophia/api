package endpoints_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wikisophia/api-arguments/arguments/argumentstest"
)

func TestGetLatest(t *testing.T) {
	expected := argumentstest.ParseSample(t, "../samples/save-request.json")
	var mistaken = expected
	mistaken.Premises = []string{"wrong", "stuff"}
	server := newServerForTests()
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
	rr := doGetArgument(newServerForTests(), 1)
	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, "text/plain; charset=utf-8", rr.Header().Get("Content-Type"))
}

func TestGetStringID(t *testing.T) {
	rr := doRequest(newServerForTests(), httptest.NewRequest("GET", "/arguments/foo", nil))
	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, "text/plain; charset=utf-8", rr.Header().Get("Content-Type"))
}
