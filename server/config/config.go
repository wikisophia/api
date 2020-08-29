package config

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// Defaults returns a Configuration with all the default options.
// This ignores environment variable values.
func Defaults() Configuration {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.ToSlash(filepath.Dir(ex))

	return Configuration{
		Server: &Server{
			Addr:                    ":8001",
			ReadHeaderTimeoutMillis: 5000,
			CorsAllowedOrigins:      []string{"*"},
			UseSSL:                  false,
			CertPath:                filepath.FromSlash(exPath + "/dev-certificates/ssl-cert.pem"),
			KeyPath:                 filepath.FromSlash(exPath + "/dev-certificates/ssl-key.pem"),
		},
		AccountsStore: &Storage{
			Type: StorageTypeMemory,
			Postgres: &Postgres{
				Database: "wikisophia_accounts",
				Host:     "localhost",
				Port:     5432,
				User:     "wikisophia_accounts_dev",
				Password: "wikisophia_accounts_dev_password",
			},
		},
		ArgumentsStore: &Storage{
			Type: StorageTypeMemory,
			Postgres: &Postgres{
				Database: "wikisophia_arguments",
				Host:     "localhost",
				Port:     5432,
				User:     "wikisophia_arguments_dev",
				Password: "wikisophia_arguments_dev_password",
			},
		},
		Hash: &Hash{
			Time:        1,
			Memory:      64 * 1024,
			Parallelism: 1,
			SaltLength:  32,
			KeyLength:   32,
		},
		JwtPrivateKeyPath: filepath.FromSlash(exPath + "/dev-certificates/jwt-private-key.pem"),
	}
}

// Configuration stores all the application config.
type Configuration struct {
	Server            *Server  `environment:"SERVER"`
	AccountsStore     *Storage `environment:"ACCOUNTS_STORE"`
	ArgumentsStore    *Storage `environment:"ARGUMENTS_STORE"`
	Hash              *Hash    `environment:"HASH"`
	JwtPrivateKeyPath string   `environment:"JWT_PRIVATE_KEY_PATH"`
}

// Server has all the config values which affect the http.Server which responds to requests.
type Server struct {
	Addr                    string   `environment:"ADDR"`
	ReadHeaderTimeoutMillis int      `environment:"READ_HEADER_TIMEOUT_MILLIS"`
	CorsAllowedOrigins      []string `environment:"CORS_ALLOWED_ORIGINS"`
	UseSSL                  bool     `environment:"USE_SSL"`
	CertPath                string   `environment:"CERT_PATH"`
	KeyPath                 string   `environment:"KEY_PATH"`
}

// Storage has all the config values related to the backend which is used to save arguments.
type Storage struct {
	Type     StorageType `environment:"TYPE"`
	Postgres *Postgres   `environment:"POSTGRES"`
}

// StorageType determines how the service stores its arguments.
type StorageType string

const (
	// StorageTypeMemory is used to save arguments unbounded, in-memory data store.
	// This is mainly intended to make development simpler, so that programmers
	// don't need to set up an actual postgres instance locally.
	StorageTypeMemory StorageType = "memory"
	// StorageTypePostgres is used to save arguments in a postgres instance.
	// If this is used, you'll need a working postgres instance to save arguments.
	StorageTypePostgres StorageType = "postgres"
)

// storageTypes returns all the valid StorageType values.
func storageTypes() []StorageType {
	return []StorageType{
		StorageTypeMemory,
		StorageTypePostgres,
	}
}

// Hash configures the hashing algorithm used to store passwords.
// This project hashes with Argon2: https://www.alexedwards.net/blog/how-to-hash-and-verify-passwords-with-argon2-in-go
type Hash struct {
	Time        uint32 `environment:"ITERATIONS"`
	Memory      uint32 `environment:"MEMORY_BYTES"`
	Parallelism uint8  `environment:"PARALLELISM"`
	SaltLength  uint8  `environment:"SALT_LENGTH"`
	KeyLength   uint32 `environment:"KEY_LENGTH"`
}

// Postgres configures the Postgres connection
type Postgres struct {
	Database string `environment:"DBNAME"`
	Host     string `environment:"HOST"`
	Port     uint16 `environment:"PORT"`
	User     string `environment:"USER"`
	Password string `environment:"PASSWORD"`
}

// ReadHeaderTimeout returns the time the server will wait for the client to send
// the HTTP headers before it just times out the request.
func (cfg *Server) ReadHeaderTimeout() time.Duration {
	return time.Duration(cfg.ReadHeaderTimeoutMillis) * time.Millisecond
}

// JwtPrivateKey returns the PrivateKey object from the file at the given path.
// This is used to sign JWTs. Panic if the file doesn't exist, can't be read, or
// didn't have a valid private key.
func (cfg *Configuration) JwtPrivateKey() *ecdsa.PrivateKey {
	data, err := ioutil.ReadFile(cfg.JwtPrivateKeyPath)
	if err != nil {
		panic("Failed to read " + prefix + "_JWT_PRIVATE_KEY_PATH file " + cfg.JwtPrivateKeyPath + ": " + err.Error())
	}

	block, _ := pem.Decode(data)
	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		panic("Couldn't parse a private key from " + cfg.JwtPrivateKeyPath + ": " + err.Error())
	}

	return key
}
