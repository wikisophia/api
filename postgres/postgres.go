package postgres

import (
	"bytes"
	"database/sql"
	"log"
	"strconv"

	// Imports the postgres driver, so that sql.Open("postgres", "blah") means something
	_ "github.com/lib/pq"
	"github.com/wikisophia/api-arguments/config"
)

// NewDB makes a connection to a postgres database.
func NewDB(cfg config.Postgres) *sql.DB {
	connStr := connectionString(cfg)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open postgres connection: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping postgres: %v", err)
	}
	return db
}

// connectionString turns the config into a string accepted by lib/pq.
// For details, see https://godoc.org/github.com/lib/pq#hdr-Connection_String_Parameters
func connectionString(cfg config.Postgres) string {
	buffer := bytes.NewBuffer(nil)

	if cfg.Host != "" {
		buffer.WriteString("host=")
		buffer.WriteString(cfg.Host)
		buffer.WriteString(" ")
	}

	if cfg.Port > 0 {
		buffer.WriteString("port=")
		buffer.WriteString(strconv.Itoa(cfg.Port))
		buffer.WriteString(" ")
	}

	if cfg.User != "" {
		buffer.WriteString("user=")
		buffer.WriteString(cfg.User)
		buffer.WriteString(" ")
	}

	if cfg.Password != "" {
		buffer.WriteString("password=")
		buffer.WriteString(cfg.Password)
		buffer.WriteString(" ")
	}

	if cfg.Database != "" {
		buffer.WriteString("dbname=")
		buffer.WriteString(cfg.Database)
		buffer.WriteString(" ")
	}

	buffer.WriteString("sslmode=disable")
	return buffer.String()
}
