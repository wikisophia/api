package postgres_test

import (
	"database/sql"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/smotes/purse"
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

	db := postgres.NewDB(config.MustParse().Storage.Postgres)
	sqlScripts, err := purse.New(filepath.Join("..", "..", "postgres", "scripts"))
	require.NoError(t, err)

	// Start with a clean slate.
	runOnce(t, sqlScripts, "destroy.sql", db)
	runOnce(t, sqlScripts, "create.sql", db)

	store := argumentsPostgres.NewPostgresStore(db)

	// Run all the same tests from the StoreTests suite.
	empty, ok := sqlScripts.Get("empty.sql")
	require.True(t, ok)
	suite.Run(t, &storetest.StoreTests{
		StoreFactory: func() arguments.Store {
			_, err := db.Query(empty)
			require.NoError(t, err)
			return store
		},
	})

	store.Close()
	db.Close()
}

func runOnce(t *testing.T, p purse.Purse, file string, db *sql.DB) {
	destroy, ok := p.Get(file)
	require.True(t, ok)
	_, err := db.Query(destroy)
	require.NoError(t, err)
}
