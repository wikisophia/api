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
	defer setEnv(t, "WKSPH_ARGS_SERVER_ADDR", "my.test.com:80")()
	defer setEnv(t, "WKSPH_ARGS_SERVER_READ_HEADER_TIMEOUT_MILLIS", "10")()
	defer setEnv(t, "WKSPH_ARGS_SERVER_CORS_ALLOWED_ORIGINS", "abc def")()
	defer setEnv(t, "WKSPH_ARGS_STORAGE_TYPE", "postgres")()
	defer setEnv(t, "WKSPH_ARGS_STORAGE_POSTGRES_DBNAME", "some-db")()
	defer setEnv(t, "WKSPH_ARGS_STORAGE_POSTGRES_HOST", "some-host")()
	defer setEnv(t, "WKSPH_ARGS_STORAGE_POSTGRES_PORT", "1234")()
	defer setEnv(t, "WKSPH_ARGS_STORAGE_POSTGRES_USER", "some-user")()
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
