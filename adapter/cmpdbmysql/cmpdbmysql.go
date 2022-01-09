package cmpdbmysql

import (
	"database/sql"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/nasjp/cmpdb"
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

func (db *mysql) Load(database *cmpdb.Database) error {
	tx, err := db.db.Begin()
	if err != nil {
		return err
	}

	if _, err = tx.Exec("SET FOREIGN_KEY_CHECKS = 0"); err != nil {
		return err
	}

	if err = db.load(tx, database); err != nil {
		return err
	}

	if _, err := tx.Exec("SET FOREIGN_KEY_CHECKS = 1"); err != nil {
		return err
	}

	return tx.Commit()
}

func (db *mysql) load(tx *sql.Tx, database *cmpdb.Database) error {
	for _, table := range database.Tables {
		if err := db.loadTable(tx, table); err != nil {
			return err
		}
	}

	return nil
}

func (db *mysql) loadTable(tx *sql.Tx, table *cmpdb.Table) error {
	for _, row := range table.Rows {
		if err := db.loadRow(tx, table.Name, row); err != nil {
			return err
		}
	}

	return nil
}

func (db *mysql) loadRow(tx *sql.Tx, tableName string, row *cmpdb.Row) error {
	_, err := tx.Exec("INSERT INTO " + tableName + fieldNames(row) + "VALUES" + fieldValues(row))
	return err
}

func fieldNames(row *cmpdb.Row) string {
	return "(" + strings.Join(row.FieldNames(), ",") + ")"
}

func fieldValues(row *cmpdb.Row) string {
	return "(" + strings.Join(row.FieldValus(), ",") + ")"
}
