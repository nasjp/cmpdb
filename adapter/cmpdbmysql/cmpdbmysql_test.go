package cmpdbmysql_test

import (
	"database/sql"
	"strings"
	"testing"

	_ "embed"

	_ "github.com/go-sql-driver/mysql"
	"github.com/nasjp/cmpdb"
	"github.com/nasjp/cmpdb/adapter/cmpdbmysql"
)

var (
	//go:embed testdata/schema.up.sql
	upSchema []byte

	//go:embed testdata/schema.down.sql
	downSchema []byte

	//go:embed testdata/diff.jsondiff
	testJSONDiff []byte
)

func TestMySQLPing(t *testing.T) {
	t.Parallel()

	db1 := getDB(t)

	setupDB(t, db1)

	if err := db1.Ping(); err != nil {
		t.Errorf("Ping() = %v, want nil", err)
	}

	comparer, err := cmpdbmysql.Load(&cmpdb.Config{
		DB:    db1,
		Bytes: testJSONDiff,
	})

	if err != nil {
		t.Fatal(err)
	}

	_ = comparer
}

func setupDB(t *testing.T, db *sql.DB) {
	t.Helper()

	for i, ddl := range splitSchema(string(upSchema)) {
		s := removeUnnecessarySchema(ddl)
		if s == "" {
			break
		}

		if _, err := db.Exec(s); err != nil {
			t.Logf("up - ddl(i: %d): '%s'", i+1, removeUnnecessarySchema(ddl))
			t.Fatal(err)
		}
	}

	t.Cleanup(func() {
		for i, ddl := range splitSchema(string(downSchema)) {
			s := removeUnnecessarySchema(ddl)
			if s == "" {
				break
			}

			if _, err := db.Exec(s); err != nil {
				t.Logf("down - ddl(i: %d): '%s'", i+1, removeUnnecessarySchema(ddl))
				t.Fatal(err)
			}
		}
	})
}

func splitSchema(schema string) []string {
	return strings.Split(schema, ";")
}

func removeUnnecessarySchema(schema string) string {
	str := strings.TrimRight(schema, "\n")
	return strings.TrimSpace(str)
}
