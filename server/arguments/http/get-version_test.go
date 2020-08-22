package http_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wikisophia/api/server/acceptancetest"
)

func TestGetVersion(t *testing.T) {
	expected := acceptancetest.ParseSample(t, samplesPath+"save-request.json")
	mistaken := expected
	mistaken.Premises = []string{"some", "bad", "version"}

	app := newApp(t, nil)
	id := app.SaveSuccessfully(t, mistaken)
	mistaken.ID = id
	mistaken.Version = 1
	expected.ID = id
	app.UpdateSuccessfully(t, expected)
	actual := app.GetVersionedSuccessfully(id, 1)
	assert.Equal(t, mistaken, actual)
}

func TestGetMissingVersion(t *testing.T) {
	arg := acceptancetest.ParseSample(t, samplesPath+"save-request.json")
	app := newApp(t, nil)
	id := app.SaveSuccessfully(t, arg)
	rr := app.Do(newGetArgumentVersion(id, 100))
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestGetStringVersion(t *testing.T) {
	rr := newApp(t, nil).Do(httptest.NewRequest("GET", "/arguments/1/version/foo", nil))
	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, "text/plain; charset=utf-8", rr.Header().Get("Content-Type"))
}

func TestGetLargeVersion(t *testing.T) {
	rr := newApp(t, nil).Do(newGetArgumentVersion(1, 65537))
	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, "text/plain; charset=utf-8", rr.Header().Get("Content-Type"))
}

func newGetArgumentVersion(id int64, version int) *http.Request {
	return httptest.NewRequest("GET", "/arguments/"+strconv.FormatInt(id, 10)+"/version/"+strconv.Itoa(version), nil)
}
