package endpoints_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVersion(t *testing.T) {
	expected := parseArgument(t, readFile(t, "../samples/save-request.json"))

	server := newServerForTests()
	id := doSaveObject(t, server, expected)
	doValidUpdate(t, server, id, []string{"some", "new", "version"})
	rr := doGetArgumentVersion(server, id, 1)
	assertSuccessfulJSON(t, rr)
	actual := parseArgument(t, rr.Body.Bytes())
	assertArgumentsMatch(t, expected, actual)
}

func TestGetMissingVersion(t *testing.T) {
	arg := parseArgument(t, readFile(t, "../samples/save-request.json"))
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
