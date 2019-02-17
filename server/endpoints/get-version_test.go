package endpoints_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wikisophia/api-arguments/server/arguments/argumentstest"
)

func TestGetVersion(t *testing.T) {
	expected := argumentstest.ParseSample(t, "../samples/save-request.json")
	mistaken := expected
	mistaken.Premises = []string{"some", "bad", "version"}

	server := newServerForTests()
	id := doSaveObject(t, server, mistaken)
	mistaken.ID = id
	expected.ID = id
	doValidUpdate(t, server, expected)
	rr := doGetArgumentVersion(server, id, 1)
	assertSuccessfulJSON(t, rr)
	actual := argumentstest.ParseJSON(t, rr.Body.Bytes())
	argumentstest.AssertArgumentsMatch(t, mistaken, actual)
}

func TestGetMissingVersion(t *testing.T) {
	arg := argumentstest.ParseSample(t, "../samples/save-request.json")
	server := newServerForTests()
	id := doSaveObject(t, server, arg)
	rr := doGetArgumentVersion(server, id, 100)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestGetStringVersion(t *testing.T) {
	rr := doRequest(newServerForTests(), httptest.NewRequest("GET", "/arguments/1/version/foo", nil))
	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, "text/plain; charset=utf-8", rr.Header().Get("Content-Type"))
}

func TestGetLargeVersion(t *testing.T) {
	rr := doGetArgumentVersion(newServerForTests(), 1, 65537)
	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, "text/plain; charset=utf-8", rr.Header().Get("Content-Type"))
}
