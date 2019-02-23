package config

import (
	"fmt"
	"log"
	"os"

	"github.com/wikisophia/go-environment-configs"
)

// MustParse wraps Parse, but prints the errors and exits rather than returning them.
func MustParse() Configuration {
	cfg, errs := Parse()
	if errs != nil {
		log.Printf("%v", errs)
		os.Exit(1)
	}
	return cfg
}

// Parse loads the app config using Viper.
// It logs all the values before returning, and panics on validation errors.
func Parse() (Configuration, error) {
	cfg := Defaults()
	errs := configs.Visit(&cfg, configs.Loader("WKSPH_ARGS"))
	log.SetOutput(os.Stdout)
	// configs.Visit(&cfg, configs.Logger("WKSPH_ARGS"))
	log.SetOutput(os.Stderr)

	errs = requirePositive(cfg.Server.ReadHeaderTimeoutMillis, "WKSPH_ARGS_SERVER_READ_HEADER_TIMEOUT_MILLIS", errs)
	errs = requirePositive(cfg.Storage.Postgres.Port, "WKSPH_ARGS_STORAGE_POSTGRES_PORT", errs)
	errs = requireValidStorageType(cfg.Storage.Type, "WKSPH_ARGS_STORAGE_TYPE", errs)
	return cfg, errs
}

// Defaults returns a Configuration with all the default options.
// This ignores environment variable values.
func Defaults() Configuration {
	return Configuration{
		Server: &Server{
			Addr:                    "localhost:8001",
			ReadHeaderTimeoutMillis: 5000,
			CorsAllowedOrigins:      []string{"*"},
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

func requirePositive(value int, prefix string, err error) error {
	if value <= 0 {
		return configs.Append(err, prefix, fmt.Errorf("must be positive. Got %d", value))
	}
	return err
}

func requireValidStorageType(value StorageType, prefix string, err error) error {
	allowedTypes := storageTypes()
	for _, storageType := range storageTypes() {
		if storageType == value {
			return err
		}
	}
	return configs.Append(err, prefix, fmt.Errorf("must be one of %v. Got %s", allowedTypes, value))
}
