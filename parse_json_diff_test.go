package cmpdb_test

import (
	"testing"

	_ "embed"

	"github.com/google/go-cmp/cmp"
	"github.com/nasjp/cmpdb"
)

//go:embed testdata/diff.jsondiff
var testJSONDiff []byte

func TestParseFromJSONDiff(t *testing.T) {
	got, err := cmpdb.ParseFromJSONDiff(testJSONDiff)
	if err != nil {
		t.Fatal(err)
	}

	want := &cmpdb.DBDiff{
		BeforeDB: &cmpdb.Database{
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
		},
		AfterDB: &cmpdb.Database{
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
									Value: &cmpdb.Value{Type: cmpdb.FieldTypeNumber, Number: 21},
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
					Name: "posts",
					Rows: []*cmpdb.Row{
						{
							Columns: []*cmpdb.Column{
								{
									Name:  "name",
									Value: &cmpdb.Value{Type: cmpdb.FieldTypeString, String: "good night"},
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
									Value: &cmpdb.Value{Type: cmpdb.FieldTypeString, String: "bye"},
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
					Name: "comments",
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
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("ParseFromJSONDiff: %s", diff)
	}

}
