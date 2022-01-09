package cmpdbmysql_test

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"testing"
	"unicode"

	_ "embed"

	"github.com/rs/zerolog"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"
)

var (
	//go:embed testdata/schema.up.sql
	upSchema []byte

	//go:embed testdata/schema.down.sql
	downSchema []byte

	//go:embed testdata/diff.jsondiff
	testJSONDiff []byte
)

const source = "root:@(localhost:3306)/%s"

var lock = &sync.Mutex{}

var dbPool = map[string]bool{
	"test_1": false,
	"test_2": false,
	"test_3": false,
}

var dbs = map[string]*sql.DB{}

// var logger = zerologadapter.New(zerolog.New(zerolog.ConsoleWriter{Out: io.Discard, NoColor: false}))

var logger = zerologadapter.New(zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false}))

func TestMain(m *testing.M) {
	if err := initTestDBs(m); err != nil {
		log.Fatalf("initDB() = %v", err)
	}
	code := m.Run()

	for _, db := range dbs {
		if err := db.Close(); err != nil {
			log.Fatalf("db.Close() = %v", err)
		}
	}

	os.Exit(code)
}

func initTestDBs(m *testing.M) error {
	grobalDB, err := initTestDB(m, "", nil)
	if err != nil {
		return err
	}

	for name := range dbPool {
		db, err := initTestDB(m, name, grobalDB)
		if err != nil {
			return err
		}

		dbs[name] = db
	}

	if err := grobalDB.Close(); err != nil {
		return err
	}

	return nil
}

func initTestDB(m *testing.M, name string, grobalDB *sql.DB) (*sql.DB, error) {
	if name != "" {
		if _, err := grobalDB.Exec("DROP DATABASE IF EXISTS " + name); err != nil {
			return nil, err
		}

		if _, err := grobalDB.Exec("CREATE DATABASE IF NOT EXISTS " + name); err != nil {
			return nil, err
		}
	}

	db, err := sql.Open("mysql", fmt.Sprintf(source, name))
	if err != nil {
		return nil, err
	}

	return sqldblogger.OpenDriver(fmt.Sprintf(source, name), db.Driver(), logger), nil
}

func getTestDB(t *testing.T) *sql.DB {
	t.Helper()

	lock.Lock()
	defer lock.Unlock()

	for {
		for name, used := range dbPool {
			if !used {
				dbPool[name] = true

				t.Cleanup(func() {
					lock.Lock()
					defer lock.Unlock()
					dbPool[name] = false
				})

				return dbs[name]
			}
		}
	}
}

func setupDB(t *testing.T, db *sql.DB) {
	t.Helper()

	execBatch(t, db, string(upSchema), "up")

	t.Cleanup(func() {
		execBatch(t, db, string(downSchema), "down")
	})
}

func execBatch(t *testing.T, db *sql.DB, schema string, msg string) {
	t.Helper()

	for i, ddl := range splitSchema(schema) {
		if _, err := db.Exec(ddl); err != nil {
			t.Logf("%s - ddl(i: %d): '%s'", msg, i+1, ddl)
			t.Fatal(err)
		}
	}
}

func splitSchema(schema string) []string {
	ddls := strings.Split(strings.TrimRightFunc(schema, unicode.IsSpace), ";")

	for i, ddl := range ddls {
		if ddl == "" {
			ddls = append(ddls[:i], ddls[i+1:]...)
		}
	}

	return ddls
}
