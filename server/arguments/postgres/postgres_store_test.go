package postgres_test

import (
	"context"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

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

	cfg := config.MustParse().ArgumentsStore.Postgres
	pool := postgres.NewPGXPool(cfg)
	emptyData, err := ioutil.ReadFile(filepath.Join(".", "scripts", "empty.sql"))
	require.NoError(t, err)
	empty := string(emptyData)
	store := argumentsPostgres.NewPostgresStore(pool)

	suite.Run(t, &storetest.StoreTests{
		StoreFactory: func() arguments.Store {
			_, err := pool.Exec(context.Background(), empty)
			require.NoError(t, err)
			return store
		},
	})

	pool.Close()
}
