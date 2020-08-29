package postgres

import (
	"bytes"
	"context"
	"database/sql"
	"log"
	"strconv"

	// Imports the postgres driver, so that sql.Open("postgres", "blah") means something

	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/wikisophia/api/server/config"
)

// NewDB makes a connection to a postgres database.
func NewDB(cfg *config.Postgres) *sql.DB {
	connStr := ParseConnectionString(cfg)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Failed to open postgres connection: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping postgres: %v", err)
	}
	return db
}

// Make a new pgx connection pool.
func NewPGXPool(cfg *config.Postgres) *pgxpool.Pool {
	pool, err := pgxpool.Connect(context.Background(), ParseConnectionString(cfg))
	if err != nil {
		log.Fatalf("Failed to open postgres connection: %v", err)
	}
	return pool
}

// MustPrepareQuery wraps db.Prepare(query), but panics on errors.
func MustPrepareQuery(db *sql.DB, query string) *sql.Stmt {
	statement, err := db.Prepare(query)
	if err != nil {
		log.Fatalf("Failed to prepare statement with query %s. Error was %v", query, err)
	}
	return statement
}

// Turn the config into a connection string.
func ParseConnectionString(cfg *config.Postgres) string {
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
