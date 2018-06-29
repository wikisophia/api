package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Storage struct {
	Type     StorageType `mapstructure:"type"`
	Postgres Postgres    `mapstructure:"postgres"`
}

func (cfg *Storage) logValues() {
	log.Printf("storage.type=%s", cfg.Type)
	if cfg.Type == StorageTypePostgres {
		cfg.Postgres.logValues()
	}
}

func (cfg *Storage) setDefaults(v *viper.Viper) {
	v.SetDefault("storage.type", "memory")
	v.SetDefault("storage.postgres.dbname", "wikisophia")
	v.SetDefault("storage.postgres.host", "localhost")
	v.SetDefault("storage.postgres.port", 5432)
	v.SetDefault("storage.postgres.user", "postgres")
	v.SetDefault("storage.postgres.password", "")
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

type StorageType string

const (
	StorageTypeMemory   StorageType = "memory"
	StorageTypePostgres StorageType = "postgres"
)

// Postgres configures the Postgres connection.
// These options come from https://godoc.org/github.com/lib/pq#hdr-Connection_String_Parameters
type Postgres struct {
	Database string `mapstructure:"dbname"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
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
