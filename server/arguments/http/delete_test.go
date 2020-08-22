package http_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wikisophia/api/server/acceptancetest"
)

func TestGetDeleted(t *testing.T) {
	app := newApp(t, nil)
	id := app.SaveSuccessfully(t, acceptancetest.ParseSample(t, samplesPath+"save-request.json"))

	rr := app.Do(newDeleteArgument(id))
	assert.Equal(t, http.StatusNoContent, rr.Code)
	assert.Equal(t, "application/json; charset=utf-8", rr.Header().Get("Content-Type"))

	rr = app.Do(newGetArgument(id))
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestDeleteUnknown(t *testing.T) {
	rr := newApp(t, nil).Do(newDeleteArgument(1))
	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, "text/plain; charset=utf-8", rr.Header().Get("Content-Type"))
}

func TestDeleteUnknownString(t *testing.T) {
	rr := newApp(t, nil).Do(httptest.NewRequest("DELETE", "/arguments/badID", nil))
	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, "text/plain; charset=utf-8", rr.Header().Get("Content-Type"))
}

func newDeleteArgument(id int64) *http.Request {
	return httptest.NewRequest("DELETE", "/arguments/"+strconv.FormatInt(id, 10), nil)
}
