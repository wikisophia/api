package config_test

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wikisophia/api-arguments/server/config"
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

	// WKSPH_ARGS_SERVER_USE_SSL determines whether the server should connect with TLS.
	assertBoolParses(t, "WKSPH_ARGS_SERVER_USE_SSL", true, func(cfg config.Configuration) bool {
		return cfg.Server.UseSSL
	})

	// WKSPH_ARGS_SERVER_KEY_PATH determines where the server should look for the key file.
	assertStringParses(t, "WKSPH_ARGS_SERVER_KEY_PATH", "/etc/ssl/certs/key.pem", func(cfg config.Configuration) string {
		return cfg.Server.KeyPath
	})

	// WKSPH_ARGS_SERVER_CERT_PATH determines where the server should look for the key file.
	assertStringParses(t, "WKSPH_ARGS_SERVER_CERT_PATH", "/etc/ssl/certs/cert.pem", func(cfg config.Configuration) string {
		return cfg.Server.CertPath
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

	// WKSPH_ARGS_HASH_ITERATIONS determines the "time" argument to argon2.Key() when hashing passwords.
	assertUInt32Parses(t, "WKSPH_ARGS_HASH_ITERATIONS", ^uint32(0), func(cfg config.Configuration) uint32 {
		return cfg.Hash.Time
	})

	// WKSPH_ARGS_HASH_MEMORY_BYTES determines the "memory" argument to argon2.Key() when hashing passwords.
	assertUInt32Parses(t, "WKSPH_ARGS_HASH_MEMORY_BYTES", ^uint32(0), func(cfg config.Configuration) uint32 {
		return cfg.Hash.Memory
	})

	// WKSPH_ARGS_HASH_MEMORY_BYTES determines the "threads" argument to argon2.Key() when hashing passwords.
	assertUInt8Parses(t, "WKSPH_ARGS_HASH_PARALLELISM", ^uint8(0), func(cfg config.Configuration) uint8 {
		return cfg.Hash.Parallelism
	})

	// WKSPH_ARGS_HASH_MEMORY_BYTES determines the length of the "salt" byte[] argument to
	// argon2.Key() when hashing passwords.
	assertUInt8Parses(t, "WKSPH_ARGS_HASH_SALT_LENGTH", ^uint8(0), func(cfg config.Configuration) uint8 {
		return cfg.Hash.SaltLength
	})

	// WKSPH_ARGS_HASH_KEY_LENGTH determines the "keyLen" argument to argon2.Key() when hashing passwords.
	assertUInt32Parses(t, "WKSPH_ARGS_HASH_KEY_LENGTH", ^uint32(0), func(cfg config.Configuration) uint32 {
		return cfg.Hash.KeyLength
	})

	// WKSPH_ARGS_JWT_PRIVATE_KEY_PATH determines which file is used as the private key to sign/verify JWTs.
	assertStringParses(t, "WKSPH_ARGS_JWT_PRIVATE_KEY_PATH", "/path-to-some-jwt-key.pem", func(cfg config.Configuration) string {
		return cfg.JwtPrivateKeyPath
	})
}

// TestLegalDefaults makes sure all the default values make a valid config object.
func TestLegalDefaults(t *testing.T) {
	_, err := config.Parse()
	assert.NoError(t, err)
}

// TestInvalidEnvironment makes sure errors show up on invalid environment variables.
func TestInvalidEnvironment(t *testing.T) {
	assertInvalid(t, "WKSPH_ARGS_SERVER_READ_HEADER_TIMEOUT_MILLIS", "foo")
	assertInvalid(t, "WKSPH_ARGS_STORAGE_POSTGRES_PORT", "foo")
	assertInvalid(t, "WKSPH_ARGS_SERVER_READ_HEADER_TIMEOUT_MILLIS", "-12")
	assertInvalid(t, "WKSPH_ARGS_SERVER_READ_HEADER_TIMEOUT_MILLIS", "0")
	assertInvalid(t, "WKSPH_ARGS_SERVER_USE_SSL", "3")
	assertInvalid(t, "WKSPH_ARGS_SERVER_USE_SSL", "notABool")
	assertInvalid(t, "WKSPH_ARGS_STORAGE_POSTGRES_PORT", "-3")
	assertInvalid(t, "WKSPH_ARGS_STORAGE_POSTGRES_PORT", "0")
	assertInvalid(t, "WKSPH_ARGS_STORAGE_TYPE", "invalid")
	assertInvalid(t, "WKSPH_ARGS_HASH_ITERATIONS", fmt.Sprintf("%d", uint64(^uint32(0))+1))
	assertInvalid(t, "WKSPH_ARGS_HASH_ITERATIONS", "-1")
	assertInvalid(t, "WKSPH_ARGS_HASH_MEMORY_BYTES", "notAnInt")
	assertInvalid(t, "WKSPH_ARGS_HASH_MEMORY_BYTES", fmt.Sprintf("%d", uint64(^uint32(0))+1))
	assertInvalid(t, "WKSPH_ARGS_HASH_MEMORY_BYTES", "-1")
	assertInvalid(t, "WKSPH_ARGS_HASH_PARALLELISM", "notAnInt")
	assertInvalid(t, "WKSPH_ARGS_HASH_PARALLELISM", fmt.Sprintf("%d", uint16(^uint8(0))+1))
	assertInvalid(t, "WKSPH_ARGS_HASH_PARALLELISM", "-1")
	assertInvalid(t, "WKSPH_ARGS_HASH_SALT_LENGTH", "notAnInt")
	assertInvalid(t, "WKSPH_ARGS_HASH_SALT_LENGTH", fmt.Sprintf("%d", uint16(^uint8(0))+1))
	assertInvalid(t, "WKSPH_ARGS_HASH_SALT_LENGTH", "-1")
	assertInvalid(t, "WKSPH_ARGS_HASH_KEY_LENGTH", "notAnInt")
	assertInvalid(t, "WKSPH_ARGS_HASH_KEY_LENGTH", fmt.Sprintf("%d", uint64(^uint32(0))+1))
	assertInvalid(t, "WKSPH_ARGS_HASH_KEY_LENGTH", "-1")
}

func TestEdgeCases(t *testing.T) {
	assertStringSliceParses(t, "WKSPH_ARGS_SERVER_CORS_ALLOWED_ORIGINS", nil, func(cfg config.Configuration) []string {
		return cfg.Server.CorsAllowedOrigins
	})
}

func assertBoolParses(t *testing.T, env string, value bool, getter func(cfg config.Configuration) bool) {
	t.Helper()
	defer setEnv(t, env, strconv.FormatBool(value))()
	cfg, errs := config.Parse()
	require.NoError(t, error(errs), "error was: \"%v\"", errs)
	assert.Equal(t, value, getter(cfg))
}

func assertStringParses(t *testing.T, env string, value string, getter func(cfg config.Configuration) string) {
	t.Helper()
	defer setEnv(t, env, value)()
	cfg, errs := config.Parse()
	require.NoError(t, error(errs), "error was: \"%v\"", errs)
	assert.Equal(t, value, getter(cfg))
}

func assertStringSliceParses(t *testing.T, env string, value []string, getter func(cfg config.Configuration) []string) {
	t.Helper()
	defer setEnv(t, env, strings.Join(value, ","))()
	cfg, errs := config.Parse()
	require.NoError(t, errs)
	assert.Equal(t, value, getter(cfg))
}

func assertIntParses(t *testing.T, env string, value int, getter func(cfg config.Configuration) int) {
	t.Helper()
	defer setEnv(t, env, strconv.Itoa(value))()
	cfg, errs := config.Parse()
	require.NoError(t, errs)
	assert.EqualValues(t, getter(cfg), value)
}

func assertUInt8Parses(t *testing.T, env string, value uint8, getter func(cfg config.Configuration) uint8) {
	t.Helper()
	defer setEnv(t, env, fmt.Sprintf("%d", value))()
	cfg, errs := config.Parse()
	require.NoError(t, errs)
	assert.EqualValues(t, getter(cfg), value)
}

func assertUInt32Parses(t *testing.T, env string, value uint32, getter func(cfg config.Configuration) uint32) {
	t.Helper()
	defer setEnv(t, env, fmt.Sprintf("%d", value))()
	cfg, errs := config.Parse()
	require.NoError(t, errs)
	assert.EqualValues(t, getter(cfg), value)
}

func assertInvalid(t *testing.T, env string, value string) {
	t.Helper()
	defer setEnv(t, env, value)()
	_, err := config.Parse()
	require.Error(t, err)
	assert.Contains(t, err.Error(), env)
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
