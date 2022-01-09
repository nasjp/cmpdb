package cmpdbmysql

import (
	"fmt"
	"strings"

	"github.com/nasjp/cmpdb"
)

func (m *mysql) Diff() (string, error) {
	currentDatabase, err := m.getCurrentDatabase()
	if err != nil {
		return "", err
	}

	_ = currentDatabase

	return "", nil
}

func (m *mysql) getCurrentDatabase() (*cmpdb.Database, error) {
	tableNames, err := m.getTableNames()
	if err != nil {
		return nil, err
	}

	table := map[string][]*Result{}

	for _, name := range tableNames {
		results, err := m.describeTable(name)
		if err != nil {
			return nil, err
		}

		table[name] = results
	}

	cmpdbTable := make([]*cmpdb.Table, 0, len(table))

	for name, results := range table {
		rows, err := m.selectRows(name, results)
		if err != nil {
			return nil, err
		}

		cmpdbTable = append(cmpdbTable, &cmpdb.Table{Name: name, Rows: rows})
	}

	return &cmpdb.Database{Tables: cmpdbTable}, nil
}

func (m *mysql) getTableNames() ([]string, error) {
	rows, err := m.db.Query("SHOW TABLES")
	if err != nil {
		return nil, err
	}

	var tableNames []string
	for rows.Next() {
		var tableName string
		if err = rows.Scan(&tableName); err != nil {
			return nil, err
		}

		tableNames = append(tableNames, tableName)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tableNames, nil
}

func (m *mysql) describeTable(name string) ([]*Result, error) {
	rows, err := m.db.Query(fmt.Sprintf("DESCRIBE %s", name))
	if err != nil {
		return nil, err
	}

	var results []*Result
	for rows.Next() {
		result := Result{}
		var tmp interface{}
		if err := rows.Scan(&result.Field, &result.Type, &result.Null, &tmp, &tmp, &tmp); err != nil {
			return nil, err
		}

		results = append(results, &result)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

type Result struct {
	Field string
	Type  string
	Null  string
}

func (m *mysql) selectRows(tableName string, results []*Result) ([]*cmpdb.Row, error) {
	rows, err := m.db.Query("SELECT * FROM " + tableName)
	if err != nil {
		return nil, err
	}

	var cmpdbRows []*cmpdb.Row

	for rows.Next() {
		// var row cmpdb.Row
		dests, columnTypes, err := genDests(results)
		if err != nil {
			return nil, err
		}

		if err := rows.Scan(dests...); err != nil {
			return nil, err
		}

		cmpdbRow, err := setRow(results, dests, columnTypes)
		if err != nil {
			return nil, err
		}

		cmpdbRows = append(cmpdbRows, cmpdbRow)
	}

	return cmpdbRows, nil
}

func genDests(results []*Result) ([]interface{}, []cmpdb.ColumnType, error) {
	var dests []interface{}
	var columnTypes []cmpdb.ColumnType
	for _, result := range results {
		v, columnType, err := fieldTypeZeroValue(result)
		if err != nil {
			return nil, nil, err
		}
		dests = append(dests, v)
		columnTypes = append(columnTypes, columnType)
	}

	return dests, columnTypes, nil
}

func fieldTypeZeroValue(result *Result) (interface{}, cmpdb.ColumnType, error) {
	switch {
	case typeMatch(result.Type, "tinyint(1)"):
		return new(bool), cmpdb.FieldTypeBoolean, nil
	case typeMatch(result.Type, "int", "tinyint", "smallint", "mediumint", "bigint"):
		return new(float64), cmpdb.FieldTypeNumber, nil
	case typeMatch(result.Type, "float", "double", "decimal", "numeric"):
		return new(float64), cmpdb.FieldTypeNumber, nil
	case typeMatch(result.Type, "varchar", "text", "char", "enum", "set", "tinytext", "mediumtext", "longtext"):
		return new(string), cmpdb.FieldTypeString, nil
	case typeMatch(result.Type, "date", "datetime", "timestamp", "time"):
		return new(string), cmpdb.FieldTypeString, nil
	default:
		return nil, 0, fmt.Errorf("unsupported type %s", result.Type)
	}
}

func typeMatch(src string, prefix ...string) bool {
	for _, t := range prefix {
		if strings.HasPrefix(src, t) {
			return true
		}
	}

	return false
}

func setRow(results []*Result, dest []interface{}, columnTypes []cmpdb.ColumnType) (*cmpdb.Row, error) {
	columns := make([]*cmpdb.Column, 0, len(dest))
	for i, v := range dest {
		if v == nil {
			columns = append(columns, &cmpdb.Column{Name: results[i].Field, Value: &cmpdb.Value{Type: cmpdb.FieldTypeNull}})
			continue
		}

		switch columnTypes[i] {
		case cmpdb.FieldTypeBoolean:
			columns = append(columns, &cmpdb.Column{Name: results[i].Field, Value: &cmpdb.Value{Boolean: *v.(*bool), Type: cmpdb.FieldTypeBoolean}})
		case cmpdb.FieldTypeNumber:
			columns = append(columns, &cmpdb.Column{Name: results[i].Field, Value: &cmpdb.Value{Number: *v.(*float64), Type: cmpdb.FieldTypeNumber}})
		case cmpdb.FieldTypeString:
			columns = append(columns, &cmpdb.Column{Name: results[i].Field, Value: &cmpdb.Value{String: *v.(*string), Type: cmpdb.FieldTypeString}})
		default:
			return nil, fmt.Errorf("unsupported type %d", columnTypes[i])
		}
	}

	return &cmpdb.Row{Columns: columns}, nil
}
