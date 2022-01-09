package cmpdb_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nasjp/cmpdb"
)

func TestDBDecoderDecode(t *testing.T) {
	want := &cmpdb.Database{
		Tables: []*cmpdb.Table{
			{
				Name: "users",
				Rows: []*cmpdb.Row{
					{
						Columns: []*cmpdb.Column{
							{
								Name:  "id",
								Value: &cmpdb.Value{Type: cmpdb.FieldTypeNumber, Number: 1},
							},
							{
								Name:  "name",
								Value: &cmpdb.Value{Type: cmpdb.FieldTypeString, String: "tom"},
							},
							{
								Name:  "age",
								Value: &cmpdb.Value{Type: cmpdb.FieldTypeNumber, Number: 20},
							},
							{
								Name:  "type",
								Value: &cmpdb.Value{Type: cmpdb.FieldTypeString, String: "manager"},
							},
						},
					},
					{
						Columns: []*cmpdb.Column{
							{
								Name:  "id",
								Value: &cmpdb.Value{Type: cmpdb.FieldTypeNumber, Number: 2},
							},
							{
								Name:  "name",
								Value: &cmpdb.Value{Type: cmpdb.FieldTypeString, String: "cindy"},
							},
							{
								Name:  "age",
								Value: &cmpdb.Value{Type: cmpdb.FieldTypeNumber, Number: 27},
							},
							{
								Name:  "type",
								Value: &cmpdb.Value{Type: cmpdb.FieldTypeString, String: "engineer"},
							},
						},
					},
				},
			},
			{
				Name: "goods",
				Rows: []*cmpdb.Row{
					{
						Columns: []*cmpdb.Column{
							{
								Name:  "name",
								Value: &cmpdb.Value{Type: cmpdb.FieldTypeString, String: "macbook"},
							},
							{
								Name:  "user_id",
								Value: &cmpdb.Value{Type: cmpdb.FieldTypeNumber, Number: 2},
							},
						},
					},
					{
						Columns: []*cmpdb.Column{
							{
								Name:  "name",
								Value: &cmpdb.Value{Type: cmpdb.FieldTypeString, String: "monitor"},
							},
							{
								Name:  "user_id",
								Value: &cmpdb.Value{Type: cmpdb.FieldTypeNumber, Number: 2},
							},
						},
					},
				},
			},
			{
				Name: "posts",
				Rows: []*cmpdb.Row{
					{
						Columns: []*cmpdb.Column{
							{
								Name:  "name",
								Value: &cmpdb.Value{Type: cmpdb.FieldTypeString, String: "wake up!"},
							},
							{
								Name:  "user_id",
								Value: &cmpdb.Value{Type: cmpdb.FieldTypeNumber, Number: 2},
							},
						},
					},
					{
						Columns: []*cmpdb.Column{
							{
								Name:  "name",
								Value: &cmpdb.Value{Type: cmpdb.FieldTypeString, String: "hello!"},
							},
							{
								Name:  "user_id",
								Value: &cmpdb.Value{Type: cmpdb.FieldTypeNumber, Number: 2},
							},
						},
					},
					{
						Columns: []*cmpdb.Column{
							{
								Name:  "name",
								Value: &cmpdb.Value{Type: cmpdb.FieldTypeString, String: "pc"},
							},
							{
								Name:  "user_id",
								Value: &cmpdb.Value{Type: cmpdb.FieldTypeNumber, Number: 1},
							},
						},
					},
				},
			},
			{
				Name: "coments",
				Rows: []*cmpdb.Row{
					{
						Columns: []*cmpdb.Column{
							{
								Name:  "name",
								Value: &cmpdb.Value{Type: cmpdb.FieldTypeString, String: "tom"},
							},
						},
					},
				},
			},
		},
	}

	got, err := cmpdb.NewDBDecoder(json.NewDecoder(bytes.NewBuffer(testBeforeJSON))).Decode()
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("DBDecoder.Decode() mismatch (-want +got):\n%s", diff)
	}
}
