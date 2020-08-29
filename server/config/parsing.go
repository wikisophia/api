package config

import (
	"log"
	"os"

	configs "github.com/wikisophia/go-environment-configs"
)

const prefix = "WKSPH_ARGS"

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
	errs := configs.LoadWithPrefix(&cfg, prefix)
	log.SetOutput(os.Stdout)
	configs.LogWithPrefix(&cfg, prefix)
	log.SetOutput(os.Stderr)

	errs = requirePositive(cfg.Server.ReadHeaderTimeoutMillis, prefix+"_SERVER_READ_HEADER_TIMEOUT_MILLIS", errs)
	errs = requirePositive(int(cfg.Storage.Postgres.Port), prefix+"_STORAGE_POSTGRES_PORT", errs)
	errs = requireValidStorageType(cfg.Storage.Type, prefix+"_STORAGE_TYPE", errs)
	return cfg, errs
}

func requirePositive(value int, prefix string, err error) error {
	return configs.Ensure(err, prefix, value > 0, "must be positive. Got %d", value)
}

func requireValidStorageType(value StorageType, prefix string, err error) error {
	allowedTypes := storageTypes()
	for _, storageType := range storageTypes() {
		if storageType == value {
			return err
		}
	}
	return configs.Ensure(err, prefix, false, "must be one of %v. Got %s", allowedTypes, value)
}
