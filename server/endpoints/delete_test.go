package endpoints_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDeleted(t *testing.T) {
	saved := parseArgument(t, readFile(t, "../samples/save-request.json"))
	server := newServerForTests()
	id := doSaveObject(t, server, saved)

	rr := doDeleteArgument(server, id)
	assert.Equal(t, http.StatusNoContent, rr.Code)
	assert.Equal(t, "application/json; charset=utf-8", rr.Header().Get("Content-Type"))

	rr = doGetArgument(server, id)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

// TODO: This test needs some database code updates. Deferring it 'til later.
// func TestDeleteUnknown(t *testing.T) {
// 	server := newServerForTests()
// 	rr := doDeleteArgument(server, 1)
// 	assert.Equal(t, http.StatusNotFound, rr.Code)
// 	assert.Equal(t, "text/plain; charset=utf-8", rr.Header().Get("Content-Type"))
// }

func TestDeleteUnknownString(t *testing.T) {
	server := newServerForTests()
	req := httptest.NewRequest("DELETE", "/arguments/badID", nil)
	rr := doRequest(server, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, "text/plain; charset=utf-8", rr.Header().Get("Content-Type"))
}
