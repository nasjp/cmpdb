package cmpdb

import (
	"database/sql"
)

type Adapter interface {
	Ping() error
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
