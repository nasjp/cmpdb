package cmpdbmysql_test

import (
	"testing"

	"github.com/nasjp/cmpdb"
	"github.com/nasjp/cmpdb/adapter/cmpdbmysql"
)

func TestMySQLLoadFixture(t *testing.T) {
	t.Parallel()

	db := getTestDB(t)

	setupDB(t, db)

	comparer, err := cmpdbmysql.Load(&cmpdb.Config{
		DB:    db,
		Bytes: testJSONDiff,
	})

	if err != nil {
		t.Fatal(err)
	}

	_ = comparer
}
