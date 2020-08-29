package postgres

import (
	"bytes"
	"context"
	"log"
	"strconv"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/wikisophia/api/server/config"
)

// Make a new pgx connection pool.
func NewPGXPool(cfg *config.Postgres) *pgxpool.Pool {
	pool, err := pgxpool.Connect(context.Background(), buildConnectionString(cfg))
	if err != nil {
		log.Fatalf("Failed to open postgres connection: %v", err)
	}
	return pool
}

// connectionString turns the config into a string accepted by pgxpool
func buildConnectionString(cfg *config.Postgres) string {
	buffer := bytes.NewBuffer(nil)

	if cfg.Host != "" {
		buffer.WriteString("host=")
		buffer.WriteString(cfg.Host)
		buffer.WriteString(" ")
	}

	if cfg.Port > 0 {
		buffer.WriteString("port=")
		buffer.WriteString(strconv.Itoa(int(cfg.Port)))
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
