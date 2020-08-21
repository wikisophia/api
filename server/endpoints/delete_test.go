package endpoints_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wikisophia/api/server/arguments/argumentstest"
)

func TestGetDeleted(t *testing.T) {
	saved := argumentstest.ParseSample(t, "../samples/save-request.json")
	server := newAppForTests(t, nil).server
	id := doSaveObject(t, server, saved)

	rr := doDeleteArgument(server, id)
	assert.Equal(t, http.StatusNoContent, rr.Code)
	assert.Equal(t, "application/json; charset=utf-8", rr.Header().Get("Content-Type"))

	rr = doGetArgument(server, id)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestDeleteUnknown(t *testing.T) {
	server := newAppForTests(t, nil).server
	rr := doDeleteArgument(server, 1)
	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, "text/plain; charset=utf-8", rr.Header().Get("Content-Type"))
}

func TestDeleteUnknownString(t *testing.T) {
	server := newAppForTests(t, nil).server
	req := httptest.NewRequest("DELETE", "/arguments/badID", nil)
	rr := doRequest(server, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, "text/plain; charset=utf-8", rr.Header().Get("Content-Type"))
}
