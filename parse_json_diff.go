package cmpdb

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/nasjp/cmpdb/jsondiff"
)

func ParseFromJSONDiff(data []byte) (*DBDiff, error) {
	beforeJSONDec, afterJSONDec, err := jsondiff.NewDecoder(data)
	if err != nil {
		return nil, err
	}

	beforeDB, err := (&dbDecoder{dec: beforeJSONDec}).decode()
	if err != nil {
		return nil, err
	}

	afterDB, err := (&dbDecoder{dec: afterJSONDec}).decode()
	if err != nil {
		return nil, err
	}

	return &DBDiff{BeforeDB: beforeDB, AfterDB: afterDB}, nil
}

type DBDiff struct {
	BeforeDB *Database
	AfterDB  *Database
}

type Database struct {
	Tables []*Table
}

type Table struct {
	Name string
	Rows []*Row
}

type Row struct {
	Columns []*Column
}

type Column struct {
	Name  string
	Value *Value
}

type ColumnType int

const (
	FieldTypeString ColumnType = iota + 1
	FieldTypeNumber
	FieldTypeBoolean
	FieldTypeNull
)

type Value struct {
	String  string
	Number  float64
	Boolean bool
	Type    ColumnType
}

type dbDecoder struct {
	dec       *json.Decoder
	currentTk json.Token
	finished  bool
}

func (d *dbDecoder) decode() (*Database, error) {
	if err := d.next(); err != nil {
		return nil, err
	}

	if d.empty() {
		return nil, errors.New("unexpected cmpdb format, db must be object")
	}

	if !d.isStart() {
		return nil, errors.New("unexpected cmpdb format, db must be object")
	}

	tables, err := d.decodeTables()
	if err != nil {
		return nil, err
	}

	return &Database{Tables: tables}, nil
}

func (d *dbDecoder) isStart() bool {
	return d.checkDelim('{')
}

func (d *dbDecoder) isEnd() bool {
	return d.checkDelim('}')
}

func (d *dbDecoder) decodeTables() ([]*Table, error) {
	tables := []*Table{}

	for {
		if err := d.next(); err != nil {
			return nil, err
		}

		if d.empty() {
			return nil, errors.New("unexpected cmpdb format")
		}

		if d.isEnd() {
			break
		}

		table, err := d.decodeTable()
		if err != nil {
			return nil, err
		}

		tables = append(tables, table)
	}

	return tables, nil
}

func (d *dbDecoder) decodeTable() (*Table, error) {
	name, err := d.decodeTableName()
	if err != nil {
		return nil, err
	}

	rows, err := d.decodeRows()
	if err != nil {
		return nil, err
	}

	return &Table{Name: name, Rows: rows}, nil
}

func (d *dbDecoder) decodeTableName() (string, error) {
	// if err := d.next(); err != nil {
	// 	return "", err
	// }

	// if d.empty() {
	// 	return "", errors.New("unexpected cmpdb format")
	// }

	name, ok := (d.currentTk).(string)
	if !ok {
		return "", errors.New("unexpected cmpdb format, table name must be string")
	}

	return name, nil
}

func (d *dbDecoder) decodeRows() ([]*Row, error) {
	if err := d.next(); err != nil {
		return nil, err
	}

	if d.empty() {
		return nil, errors.New("unexpected cmpdb format")
	}

	if !d.isRowsStart() {
		return nil, errors.New("unexpected cmpdb format, table rows must be array")
	}

	rows := []*Row{}

	for {
		if err := d.next(); err != nil {
			return nil, err
		}

		if d.empty() {
			return nil, errors.New("unexpected cmpdb format")
		}

		if d.isRowsEnd() {
			break
		}

		row, err := d.decodeRow()
		if err != nil {
			return nil, err
		}

		rows = append(rows, row)
	}

	if !d.isRowsEnd() {
		return nil, errors.New("unexpected cmpdb format, table rows must be array")
	}

	return rows, nil
}

func (d *dbDecoder) isRowsStart() bool {
	return d.checkDelim('[')
}

func (d *dbDecoder) isRowsEnd() bool {
	return d.checkDelim(']')
}

func (d *dbDecoder) decodeRow() (*Row, error) {
	// if err := d.next(); err != nil {
	// 	return nil, err
	// }

	// if d.empty() {
	// 	return nil, errors.New("unexpected cmpdb format")
	// }

	if !d.isRowStart() {
		return nil, errors.New("unexpected cmpdb format, row must be object")
	}

	columns := []*Column{}

	for {

		if err := d.next(); err != nil {
			return nil, err
		}

		if d.empty() {
			return nil, errors.New("unexpected cmpdb format")
		}

		if d.isRowEnd() {
			break
		}

		column, err := d.decodeColumn()
		if err != nil {
			return nil, err
		}

		columns = append(columns, column)
	}

	if !d.isRowEnd() {
		return nil, errors.New("unexpected cmpdb format, row must be object")
	}

	return &Row{Columns: columns}, nil
}

func (d *dbDecoder) isRowStart() bool {
	return d.checkDelim('{')
}

func (d *dbDecoder) isRowEnd() bool {
	return d.checkDelim('}')
}

func (d *dbDecoder) decodeColumn() (*Column, error) {
	name, err := d.decodeColumnName()
	if err != nil {
		return nil, err
	}

	value, err := d.decodeColumnValue()
	if err != nil {
		return nil, err
	}

	return &Column{Name: name, Value: value}, nil
}

func (d *dbDecoder) decodeColumnName() (string, error) {
	// if err := d.next(); err != nil {
	// 	return "", err
	// }

	// if d.empty() {
	// 	return "", errors.New("unexpected cmpdb format")
	// }

	name, ok := (d.currentTk).(string)
	if !ok {
		return "", errors.New("unexpected cmpdb format, column name must be string")
	}

	return name, nil
}

func (d *dbDecoder) decodeColumnValue() (*Value, error) {
	if err := d.next(); err != nil {
		return nil, err
	}

	if d.empty() {
		return nil, errors.New("unexpected cmpdb format")
	}

	switch v := (d.currentTk).(type) {
	case string:
		return &Value{String: v, Type: FieldTypeString}, nil
	case float64:
		return &Value{Number: v, Type: FieldTypeNumber}, nil
	case int:
		return &Value{Number: float64(v), Type: FieldTypeNumber}, nil
	default:
		return nil, errors.New("unexpected cmpdb format, value must be string or number")
	}
}

func (d *dbDecoder) empty() bool {
	return d.currentTk == nil
}

func (d *dbDecoder) next() error {
	token, err := d.dec.Token()
	if err != nil {
		if err != io.EOF {
			return err
		}

		d.finished = true

		return nil
	}

	d.currentTk = token

	return nil
}

func (d *dbDecoder) checkDelim(delim json.Delim) bool {
	tk, ok := (d.currentTk).(json.Delim)
	if ok && tk == delim {
		return true
	}

	return false
}
