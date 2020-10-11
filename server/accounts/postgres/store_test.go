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
	"github.com/wikisophia/api/server/accounts"
	accountsPostgres "github.com/wikisophia/api/server/accounts/postgres"
	"github.com/wikisophia/api/server/accounts/storetest"
	"github.com/wikisophia/api/server/config"
	"github.com/wikisophia/api/server/passwords"
	"github.com/wikisophia/api/server/postgres"
)

var hasDatabase = flag.Bool("database", false, "run database integration tests")

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

// TestPostgresStore makes sure store is consistent with the StoreTests suite.
// This helps verify:
//    1. The postgres Store, which is used in prod to keep account info.
//    2.The StoreTests suite, which is reused to test the inMemoryStore as well.
func TestPostgresStore(t *testing.T) {
	if !*hasDatabase {
		return
	}

	cfg := config.MustParse()
	pool := postgres.NewPGXPool(cfg.AccountsStore.Postgres)
	emptyData, err := ioutil.ReadFile(filepath.Join(".", "scripts", "empty.sql"))
	require.NoError(t, err)
	empty := string(emptyData)
	store := accountsPostgres.NewPostgresStore(pool, passwords.NewHasher(*cfg.Hash))

	suite.Run(t, &storetest.StoreTests{
		StoreFactory: func() accounts.Store {
			_, err := pool.Exec(context.Background(), empty)
			require.NoError(t, err)
			return store
		},
	})

	pool.Close()
}
