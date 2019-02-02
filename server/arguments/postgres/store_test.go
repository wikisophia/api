package postgres_test

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/wikisophia/api-arguments/server/arguments"
	"github.com/wikisophia/api-arguments/server/arguments/argumentstest"

	"github.com/smotes/purse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	argumentsInPostgres "github.com/wikisophia/api-arguments/server/arguments/postgres"
	"github.com/wikisophia/api-arguments/server/config"
	"github.com/wikisophia/api-arguments/server/postgres"
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

	// These tests will be slow... so do as much as we can up front to save time
	sqlScripts, err := purse.New(filepath.Join("..", "..", "postgres", "scripts"))
	if !assert.NoError(t, err) {
		return
	}
	contents, ok := sqlScripts.Get("clear.sql")
	if !assert.True(t, ok) {
		return
	}

	db := postgres.NewDB(config.MustParse().Storage.Postgres)
	if _, err = db.Query(contents); !assert.NoError(t, err) {
		return
	}

	store := argumentsInPostgres.NewStore(db)
	// Run all the same tests from the StoreTests suite.
	suite.Run(t, &argumentstest.StoreTests{
		StoreFactory: func() arguments.Store {
			if _, err := db.Query(contents); !assert.NoError(t, err) {
				t.FailNow()
			}
			return store
		},
	})

	store.Close()
	db.Close()
}
