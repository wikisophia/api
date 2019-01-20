package postgres_test

import (
	"database/sql"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/wikisophia/api-arguments/arguments/argumentstest"

	"github.com/smotes/purse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	argumentsInPostgres "github.com/wikisophia/api-arguments/arguments/postgres"
	"github.com/wikisophia/api-arguments/config"
	"github.com/wikisophia/api-arguments/postgres"
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
	if _, err = db.Query(contents); assert.NoError(t, err) {
		return
	}

	store := argumentsInPostgres.NewStore(db)
	// Run all the same tests from the StoreTests suite.
	suite.Run(t, &DatabaseTests{
		StoreTests: &argumentstest.StoreTests{
			Store: store,
		},
		db:            db,
		queryContents: contents,
	})

	store.Close()
	db.Close()
}

type DatabaseTests struct {
	*argumentstest.StoreTests

	db            *sql.DB
	queryContents string
}

func (suite *DatabaseTests) SetupTest() {
	if suite.db != nil {
		if _, err := suite.db.Query(suite.queryContents); !assert.NoError(suite.T(), err) {
			os.Exit(1)
		}
	}
}
