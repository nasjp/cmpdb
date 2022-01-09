package cmpdbmysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type mysql struct {
	db *sql.DB
}

func New(db *sql.DB) *mysql {
	return &mysql{db}
}

func (db *mysql) Ping() error {
	return db.db.Ping()
}
