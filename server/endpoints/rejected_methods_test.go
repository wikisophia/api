package endpoints_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostWithID(t *testing.T) {
	req := httptest.NewRequest("POST", "/arguments/1", nil)
	rr := doRequest(newAppForTests(testServerConfig{}).server, req)
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func TestPostVersion(t *testing.T) {
	req := httptest.NewRequest("POST", "/arguments/1/version/1", nil)
	rr := doRequest(newAppForTests(testServerConfig{}).server, req)
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}
