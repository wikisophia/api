package arguments_test

import (
	"database/sql"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/wikisophia/api-arguments/server/endpoints"

	"github.com/smotes/purse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/wikisophia/api-arguments/server/arguments"
	"github.com/wikisophia/api-arguments/server/arguments/argumentstest"
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

	db := postgres.NewDB(config.MustParse().Storage.Postgres)
	sqlScripts, err := purse.New(filepath.Join(".", "scripts"))
	if !assert.NoError(t, err) {
		return
	}

	// Start with a clean slate.
	if !runOnce(t, sqlScripts, "destroy.sql", db) {
		return
	}
	if !runOnce(t, sqlScripts, "create.sql", db) {
		return
	}

	store := arguments.NewPostgresStore(db)

	// Run all the same tests from the StoreTests suite.
	empty, ok := sqlScripts.Get("empty.sql")
	if !assert.True(t, ok) {
		return
	}
	suite.Run(t, &argumentstest.StoreTests{
		StoreFactory: func() endpoints.Store {
			if _, err := db.Query(empty); !assert.NoError(t, err) {
				t.FailNow()
			}
			return store
		},
	})

	store.Close()
	db.Close()
}

func runOnce(t *testing.T, p purse.Purse, file string, db *sql.DB) bool {
	destroy, ok := p.Get(file)
	if !assert.True(t, ok) {
		return false
	}
	if _, err := db.Query(destroy); !assert.NoError(t, err) {
		return false
	}
	return true
}
