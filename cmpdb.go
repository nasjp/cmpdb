package cmpdb

import (
	"database/sql"
)

type Adapter interface {
	LoadFixture(database *Database) error
}

type Config struct {
	DB    *sql.DB
	Bytes []byte
}

type Comparer struct {
	Adapter Adapter
	Bytes   []byte
	DBDiff  *DBDiff
}
