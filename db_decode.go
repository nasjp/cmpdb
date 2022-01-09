package cmpdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

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

type DBDecoder struct {
	dec       *json.Decoder
	currentTk json.Token
	finished  bool
}

func NewDBDecoder(dec *json.Decoder) *DBDecoder {
	return &DBDecoder{dec: dec}
}

func (d *DBDecoder) Decode() (*Database, error) {
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

func (d *DBDecoder) isStart() bool {
	return d.checkDelim('{')
}

func (d *DBDecoder) isEnd() bool {
	return d.checkDelim('}')
}

func (d *DBDecoder) decodeTables() ([]*Table, error) {
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

func (d *DBDecoder) decodeTable() (*Table, error) {
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

func (d *DBDecoder) decodeTableName() (string, error) {
	name, ok := (d.currentTk).(string)
	if !ok {
		return "", errors.New("unexpected cmpdb format, table name must be string")
	}

	return name, nil
}

func (d *DBDecoder) decodeRows() ([]*Row, error) {
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

func (d *DBDecoder) isRowsStart() bool {
	return d.checkDelim('[')
}

func (d *DBDecoder) isRowsEnd() bool {
	return d.checkDelim(']')
}

func (d *DBDecoder) decodeRow() (*Row, error) {
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

func (d *DBDecoder) isRowStart() bool {
	return d.checkDelim('{')
}

func (d *DBDecoder) isRowEnd() bool {
	return d.checkDelim('}')
}

func (d *DBDecoder) decodeColumn() (*Column, error) {
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

func (d *DBDecoder) decodeColumnName() (string, error) {
	name, ok := (d.currentTk).(string)
	if !ok {
		return "", errors.New("unexpected cmpdb format, column name must be string")
	}

	return name, nil
}

func (d *DBDecoder) decodeColumnValue() (*Value, error) {
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

func (d *DBDecoder) empty() bool {
	return d.currentTk == nil
}

func (d *DBDecoder) next() error {
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

func (d *DBDecoder) checkDelim(delim json.Delim) bool {
	tk, ok := (d.currentTk).(json.Delim)
	if ok && tk == delim {
		return true
	}

	return false
}

func (r *Row) FieldNames() []string {
	names := make([]string, 0, len(r.Columns))

	for _, column := range r.Columns {
		names = append(names, column.Name)
	}

	return names
}

func (r *Row) FieldValus() []string {
	fields := make([]string, 0, len(r.Columns))

	for _, column := range r.Columns {
		fields = append(fields, column.ValueString())
	}

	return fields
}

func (c *Column) ValueString() string {
	switch c.Value.Type {
	case FieldTypeString:
		return `"` + c.Value.String + `"`
	case FieldTypeNumber:
		return fmt.Sprintf("%f", c.Value.Number)
	case FieldTypeBoolean:
		return fmt.Sprintf("%t", c.Value.Boolean)
	case FieldTypeNull:
		return "null"
	}

	// unreachable
	return ""
}
