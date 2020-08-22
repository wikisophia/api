package acceptancetest

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertBadRequest(t *testing.T, method, path, body string) {
	t.Helper()
	a := NewApp(t, nil)
	rr := a.Do(httptest.NewRequest(method, path, strings.NewReader(body)))
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func AssertMethodNotAllowed(t *testing.T, method, path string) {
	t.Helper()
	a := NewApp(t, nil)
	rr := a.Do(httptest.NewRequest(method, path, nil))
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}
