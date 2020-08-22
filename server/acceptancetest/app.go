package acceptancetest

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	accountsMemory "github.com/wikisophia/api/server/accounts/memory"
	argumentsMemory "github.com/wikisophia/api/server/arguments/memory"
	wikisophiaHttp "github.com/wikisophia/api/server/http"
)

// NewApp returns a bundle of utils useful for acceptance testing the app.
func NewApp(t *testing.T, cfg *AppConfig) *App {
	if cfg == nil {
		cfg = &AppConfig{
			EmailerSucceeds: true,
		}
	}
	emailer := &Emailer{
		shouldSucceed: cfg.EmailerSucceeds,
	}
	server := wikisophiaHttp.NewServer(newKeyForTests(t), wikisophiaHttp.ServerDependencies{
		AccountsStore:  accountsMemory.NewMemoryStore(),
		ArgumentsStore: argumentsMemory.NewMemoryStore(),
		Emailer:        emailer,
	})
	return &App{
		t:       t,
		server:  server,
		Emailer: emailer,
	}
}

type App struct {
	t       *testing.T
	server  *wikisophiaHttp.Server
	Emailer *Emailer
}

type AppConfig struct {
	EmailerSucceeds bool
}

func (a *App) Do(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.server.Handle(rr, req)
	return rr
}

func (a *App) AssertBadRequest(method, path, body string) {
	a.t.Helper()
	rr := a.Do(httptest.NewRequest(method, path, strings.NewReader(body)))
	assert.Equal(a.t, http.StatusBadRequest, rr.Code)
	assert.Equal(a.t, "text/plain; charset=utf-8", rr.Header().Get("Content-Type"))
}

func (a *App) AssertNotFound(method, path string) {
	a.t.Helper()
	rr := a.Do(httptest.NewRequest(method, path, nil))
	assert.Equal(a.t, http.StatusNotFound, rr.Code)
}

func newKeyForTests(t *testing.T) *ecdsa.PrivateKey {
	data := `-----BEGIN EC PRIVATE KEY-----
MIGkAgEBBDAAT0uP+yvZG/miupKFfwmEetHH71/aBcwZLcThBleGSp7pj3Y3lsAq
dGj/LESq07ygBwYFK4EEACKhZANiAATxCwIvHjShDKaLqv9sHDLelsjVlC5YoRnZ
T7uz6EGYcGfvnuIkD12uolnsuqvgzybhxFw2/311B1t7eXNJEBh6VdqkC4k8DhhH
BMLqx6d62CqjE3PIUXzJR9mtgNo1PrA=
-----END EC PRIVATE KEY-----
`
	block, _ := pem.Decode([]byte(data))
	key, err := x509.ParseECPrivateKey(block.Bytes)
	require.NoError(t, err, "Error generating private key for tests")
	return key
}
