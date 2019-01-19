package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wikisophia/api-arguments/config"
)

func TestLegalDefaults(t *testing.T) {
	_, errs := config.ParseConfig()
	assert.Len(t, errs, 0)
}

func TestEnvironmentOverrides(t *testing.T) {
	// WKSPH_ARGS_SERVER_ADDR determines which host/port the server attaches to.
	defer setEnv(t, "WKSPH_ARGS_SERVER_ADDR", "my.test.com:80")()

	// WKSPH_ARGS_SERVER_READ_HEADER_TIMEOUT_MILLIS determines the number of milliseconds
	// the server will wait for a client to send the request headers before timing out.
	defer setEnv(t, "WKSPH_ARGS_SERVER_READ_HEADER_TIMEOUT_MILLIS", "10")()

	// WKSPH_ARGS_SERVER_CORS_ALLOWED_ORIGINS configures which domains can call us
	// with a CORS request. For more info, see https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS
	defer setEnv(t, "WKSPH_ARGS_SERVER_CORS_ALLOWED_ORIGINS", "abc def")()

	// WKSPH_ARGS_STORAGE_TYPE determines how the service stores data. Valid options are "memory" or "postgres".
	defer setEnv(t, "WKSPH_ARGS_STORAGE_TYPE", "postgres")()

	// WKSPH_ARGS_STORAGE_POSTGRES_DBNAME determines which database inside postgres is used.
	// If WKSPH_ARGS_STORAGE_TYPE is "memory", this is ignored.
	defer setEnv(t, "WKSPH_ARGS_STORAGE_POSTGRES_DBNAME", "some-db")()

	// WKSPH_ARGS_STORAGE_POSTGRES_HOST determines which hostname the service should use when connecting to postgres.
	// If WKSPH_ARGS_STORAGE_TYPE is "memory", this is ignored.
	defer setEnv(t, "WKSPH_ARGS_STORAGE_POSTGRES_HOST", "some-host")()

	// WKSPH_ARGS_STORAGE_POSTGRES_PORT determines which port the service should use when connecting to postgres.
	// If WKSPH_ARGS_STORAGE_TYPE is "memory", this is ignored.
	defer setEnv(t, "WKSPH_ARGS_STORAGE_POSTGRES_PORT", "1234")()

	// WKSPH_ARGS_STORAGE_POSTGRES_USER determines which username the service should use when connecting to postgres.
	// If WKSPH_ARGS_STORAGE_TYPE is "memory", this is ignored.
	defer setEnv(t, "WKSPH_ARGS_STORAGE_POSTGRES_USER", "some-user")()

	// WKSPH_ARGS_STORAGE_POSTGRES_PASSWORD determines which password the service should use when connecting to postgres.
	// If WKSPH_ARGS_STORAGE_TYPE is "memory", this is ignored.
	defer setEnv(t, "WKSPH_ARGS_STORAGE_POSTGRES_PASSWORD", "some-password")()

	cfg, errs := config.ParseConfig()

	assert.Len(t, errs, 0)
	assert.Equal(t, cfg.Server.Addr, "my.test.com:80")
	assert.Equal(t, cfg.Server.ReadHeaderTimeoutMillis, 10)
	assert.Equal(t, cfg.Server.CorsAllowedOrigins, []string{"abc", "def"})
	assert.Equal(t, cfg.Storage.Type, config.StorageTypePostgres)
	assert.Equal(t, cfg.Storage.Postgres.Database, "some-db")
	assert.Equal(t, cfg.Storage.Postgres.Port, 1234)
	assert.Equal(t, cfg.Storage.Postgres.User, "some-user")
	assert.Equal(t, cfg.Storage.Postgres.Password, "some-password")
}

// setEnv acts as a wrapper around os.Setenv, returning a function that resets the environment
// back to its original value. This prevents tests from setting environment variables as a side-effect.
func setEnv(t *testing.T, key string, val string) func() {
	orig, set := os.LookupEnv(key)
	err := os.Setenv(key, val)
	if !assert.NoError(t, err) {
		return func() {}
	}
	if set {
		return func() {
			os.Setenv(key, orig)
		}
	} else {
		return func() {
			os.Unsetenv(key)
		}
	}
}
