package config

import (
	"fmt"
	"log"
)

func (cfg *Storage) logValues() {
	log.Printf("storage.type=%s", cfg.Type)
	if cfg.Type == StorageTypePostgres {
		cfg.Postgres.logValues()
	}
}

func (cfg *Storage) validate() []error {
	switch cfg.Type {
	case StorageTypeMemory:
		return nil
	case StorageTypePostgres:
		return cfg.Postgres.validate()
	default:
		return []error{fmt.Errorf(`storage.type has unrecognized value: %s. This must be either "memory" or "postgres"`, cfg.Type)}
	}
}

func parseStorageType(value string) (StorageType, error) {
	switch value {
	case "memory":
		return StorageTypeMemory, nil
	case "postgres":
		return StorageTypePostgres, nil
	default:
		return StorageTypeMemory, fmt.Errorf("%s must be one of \"memory\" or \"postgres\"", value)
	}
}

func (cfg *Postgres) logValues() {
	log.Printf("postgres.dbname=%s", cfg.Database)
	log.Printf("postgres.host=%s", cfg.Host)
	log.Printf("postgres.port=%d", cfg.Port)
	log.Printf("postgres.user=%s", cfg.User)
	// Don't log the password, for security reasons
	// log.Printf("postgres.password=%s", cfg.Password)
}

func (cfg *Postgres) validate() []error {
	return nil
}
