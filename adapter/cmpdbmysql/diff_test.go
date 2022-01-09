package cmpdbmysql_test

import (
	"testing"

	"github.com/nasjp/cmpdb/adapter/cmpdbmysql"
)

func TestMySQLDiff(t *testing.T) {
	t.Parallel()

	db := getTestDB(t)

	setupDB(t, db)

	execBatch(t, db, "INSERT INTO users (name, age, type) VALUES ('user1', 10, 'customer');", "insert")

	mySQL := cmpdbmysql.New(db)

	if _, err := mySQL.Diff(); err != nil {
		t.Fatal(err)
	}
}
