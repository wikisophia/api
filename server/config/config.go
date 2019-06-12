package config

import (
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
			Addr: ":8001",
			ReadHeaderTimeoutMillis: 5000,
			CorsAllowedOrigins:      []string{"*"},
			UseSSL:                  false,
			CertPath:                filepath.FromSlash(exPath + "/certificates/cert.pem"),
			KeyPath:                 filepath.FromSlash(exPath + "/certificates/key.pem"),
		},
		Storage: &Storage{
			Type: StorageTypeMemory,
			Postgres: &Postgres{
				Database: "wikisophia",
				Host:     "localhost",
				Port:     5432,
				User:     "postgres",
				Password: "",
			},
		},
	}
}

// Configuration stores all the application config.
type Configuration struct {
	Server  *Server  `environment:"SERVER"`
	Storage *Storage `environment:"STORAGE"`
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

// Postgres configures the Postgres connection.
// These options come from https://godoc.org/github.com/lib/pq#hdr-Connection_String_Parameters
type Postgres struct {
	Database string `environment:"DBNAME"`
	Host     string `environment:"HOST"`
	Port     int    `environment:"PORT"`
	User     string `environment:"USER"`
	Password string `environment:"PASSWORD"`
}

// ReadHeaderTimeout returns the time the server will wait for the client to send
// the HTTP headers before it just times out the request.
func (cfg *Server) ReadHeaderTimeout() time.Duration {
	return time.Duration(cfg.ReadHeaderTimeoutMillis) * time.Millisecond
}
