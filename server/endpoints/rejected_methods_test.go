package endpoints_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCollection(t *testing.T) {
	rr := doRequest(newServerForTests(), httptest.NewRequest("GET", "/arguments", nil))
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestPostWithID(t *testing.T) {
	req := httptest.NewRequest("POST", "/arguments/1", nil)
	rr := doRequest(newServerForTests(), req)
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func TestPostVersion(t *testing.T) {
	req := httptest.NewRequest("POST", "/arguments/1/version/1", nil)
	rr := doRequest(newServerForTests(), req)
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}
