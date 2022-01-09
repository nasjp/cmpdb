package cmpdbmysql_test

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestMySQLPing(t *testing.T) {
	t.Parallel()

	db1 := getDB(t)

	if err := db1.Ping(); err != nil {
		t.Errorf("Ping() = %v, want nil", err)
	}
}
