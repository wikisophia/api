package config_test

import (
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wikisophia/api-arguments/config"
)

func TestEnvironmentOverrides(t *testing.T) {
	// WKSPH_ARGS_SERVER_ADDR determines which host/port the server attaches to.
	assertStringParses(t, "WKSPH_ARGS_SERVER_ADDR", "my.test.com:80", func(cfg config.Configuration) string {
		return cfg.Server.Addr
	})

	// WKSPH_ARGS_SERVER_READ_HEADER_TIMEOUT_MILLIS determines the number of milliseconds
	// the server will wait for a client to send the request headers before timing out.
	assertIntParses(t, "WKSPH_ARGS_SERVER_READ_HEADER_TIMEOUT_MILLIS", 10, func(cfg config.Configuration) int {
		return cfg.Server.ReadHeaderTimeoutMillis
	})

	// WKSPH_ARGS_SERVER_CORS_ALLOWED_ORIGINS configures which domains can call us
	// with a CORS request. For more info, see https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS
	assertStringSliceParses(t, "WKSPH_ARGS_SERVER_CORS_ALLOWED_ORIGINS", []string{"abc", "def"}, func(cfg config.Configuration) []string {
		return cfg.Server.CorsAllowedOrigins
	})

	// WKSPH_ARGS_STORAGE_TYPE determines how the service stores data. Valid options are "memory" or "postgres".
	assertStringParses(t, "WKSPH_ARGS_STORAGE_TYPE", "postgres", func(cfg config.Configuration) string {
		return string(cfg.Storage.Type)
	})

	// WKSPH_ARGS_STORAGE_POSTGRES_DBNAME determines which database inside postgres is used.
	// If WKSPH_ARGS_STORAGE_TYPE is "memory", this is ignored.
	assertStringParses(t, "WKSPH_ARGS_STORAGE_POSTGRES_DBNAME", "some-db", func(cfg config.Configuration) string {
		return cfg.Storage.Postgres.Database
	})

	// WKSPH_ARGS_STORAGE_POSTGRES_HOST determines which hostname the service should use when connecting to postgres.
	// If WKSPH_ARGS_STORAGE_TYPE is "memory", this is ignored.
	assertStringParses(t, "WKSPH_ARGS_STORAGE_POSTGRES_HOST", "some-host", func(cfg config.Configuration) string {
		return cfg.Storage.Postgres.Host
	})

	// WKSPH_ARGS_STORAGE_POSTGRES_PORT determines which port the service should use when connecting to postgres.
	// If WKSPH_ARGS_STORAGE_TYPE is "memory", this is ignored.
	assertIntParses(t, "WKSPH_ARGS_STORAGE_POSTGRES_PORT", 1234, func(cfg config.Configuration) int {
		return cfg.Storage.Postgres.Port
	})

	// WKSPH_ARGS_STORAGE_POSTGRES_USER determines which username the service should use when connecting to postgres.
	// If WKSPH_ARGS_STORAGE_TYPE is "memory", this is ignored.
	assertStringParses(t, "WKSPH_ARGS_STORAGE_POSTGRES_USER", "some-user", func(cfg config.Configuration) string {
		return cfg.Storage.Postgres.User
	})

	// WKSPH_ARGS_STORAGE_POSTGRES_PASSWORD determines which password the service should use when connecting to postgres.
	// If WKSPH_ARGS_STORAGE_TYPE is "memory", this is ignored.
	assertStringParses(t, "WKSPH_ARGS_STORAGE_POSTGRES_PASSWORD", "some-password", func(cfg config.Configuration) string {
		return cfg.Storage.Postgres.Password
	})
}

// TestLegalDefaults makes sure all the default values make a valid config object.
func TestLegalDefaults(t *testing.T) {
	_, errs := config.Parse()
	assert.Len(t, errs, 0)
}

// TestInvalidEnvironment makes sure
func TestInvalidEnvironment(t *testing.T) {
	assertInvalid(t, "WKSPH_ARGS_SERVER_READ_HEADER_TIMEOUT_MILLIS", "foo")
	assertInvalid(t, "WKSPH_ARGS_STORAGE_POSTGRES_PORT", "foo")
	assertInvalid(t, "WKSPH_ARGS_SERVER_READ_HEADER_TIMEOUT_MILLIS", "-12")
	assertInvalid(t, "WKSPH_ARGS_STORAGE_POSTGRES_PORT", "-3")
}

func TestEdgeCases(t *testing.T) {
	assertStringSliceParses(t, "WKSPH_ARGS_SERVER_CORS_ALLOWED_ORIGINS", nil, func(cfg config.Configuration) []string {
		return cfg.Server.CorsAllowedOrigins
	})
	assertInvalid(t, "WKSPH_ARGS_SERVER_READ_HEADER_TIMEOUT_MILLIS", "0")
	assertInvalid(t, "WKSPH_ARGS_STORAGE_POSTGRES_PORT", "0")
}

func assertStringParses(t *testing.T, env string, value string, getter func(cfg config.Configuration) string) {
	t.Helper()
	defer setEnv(t, env, value)()
	cfg, errs := config.Parse()
	if !assert.Len(t, errs, 0) {
		return
	}
	assert.Equal(t, value, getter(cfg))
}

func assertStringSliceParses(t *testing.T, env string, value []string, getter func(cfg config.Configuration) []string) {
	t.Helper()
	defer setEnv(t, env, strings.Join(value, ","))()
	cfg, errs := config.Parse()
	if !assert.Len(t, errs, 0) {
		return
	}
	assert.Equal(t, value, getter(cfg))
}

func assertIntParses(t *testing.T, env string, value int, getter func(cfg config.Configuration) int) {
	t.Helper()
	defer setEnv(t, env, strconv.Itoa(value))()
	cfg, errs := config.Parse()
	if !assert.Len(t, errs, 0) {
		return
	}
	assert.EqualValues(t, getter(cfg), value)
}

func assertInvalid(t *testing.T, env string, value string) {
	t.Helper()
	defer setEnv(t, env, value)()
	_, errs := config.Parse()
	assert.Len(t, errs, 1)
	assert.True(t, strings.HasPrefix(errs[0].Error(), env), "error message \"%v\" must start with environment variable: \"%s\"", errs[0], env)
}

// setEnv acts as a wrapper around os.Setenv, returning a function that resets the environment
// back to its original value. This prevents tests from setting environment variables as a side-effect.
func setEnv(t *testing.T, key string, val string) func() {
	t.Helper()
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
