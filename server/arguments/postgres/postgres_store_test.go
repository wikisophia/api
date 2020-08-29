package postgres_test

import (
	"context"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/wikisophia/api/server/arguments"
	argumentsPostgres "github.com/wikisophia/api/server/arguments/postgres"
	"github.com/wikisophia/api/server/arguments/storetest"
	"github.com/wikisophia/api/server/config"
	"github.com/wikisophia/api/server/postgres"
)

var hasDatabase = flag.Bool("database", false, "run database integration tests")

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func TestArgumentStorageIntegration(t *testing.T) {
	// Only run tests which rely on the database if the database flag is present
	if !*hasDatabase {
		return
	}

	pool := postgres.NewPGXPool(config.MustParse().Storage.Postgres)
	create := mustReadScript(t, "create.sql")
	destroy := mustReadScript(t, "destroy.sql")
	empty := mustReadScript(t, "empty.sql")

	// Start with a clean slate.
	mustRun(t, destroy, pool)
	mustRun(t, create, pool)
	store := argumentsPostgres.NewPostgresStore(pool)

	// Run all the same tests from the StoreTests suite.
	suite.Run(t, &storetest.StoreTests{
		StoreFactory: func() arguments.Store {
			mustRun(t, empty, pool)
			return store
		},
	})

	store.Close()
	pool.Close()
}

func mustReadScript(t *testing.T, filename string) string {
	data, err := ioutil.ReadFile(filepath.Join("..", "..", "postgres", "scripts", filename))
	require.NoError(t, err)
	return string(data)
}

func mustRun(t *testing.T, commands string, pool *pgxpool.Pool) {
	_, err := pool.Exec(context.Background(), commands)
	require.NoError(t, err)
}
