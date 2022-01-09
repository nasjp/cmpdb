package cmpdbmysql_test

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"testing"
)

const source = "root:@(localhost:3306)/%s"

var lock = &sync.Mutex{}

var dbPool = map[string]bool{
	"test_1": false,
	"test_2": false,
	"test_3": false,
}

var dbs = map[string]*sql.DB{}

func TestMain(m *testing.M) {
	if err := initTestDBs(m); err != nil {
		log.Fatalf("initDB() = %v", err)
	}
	code := m.Run()

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

	return db, nil
}

func getDB(t *testing.T) *sql.DB {
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
