package cmpdb

import (
	"database/sql"

	"github.com/nasjp/cmpdb/adapter"
)

type Config struct {
	DB    *sql.DB
	Bytes []byte
}

type Comparer struct {
	Adapter adapter.Adapter
}
