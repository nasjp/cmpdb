package cmpdb_test

import (
	_ "embed"
)

var (

	//go:embed testdata/diff.jsondiff
	testJSONDiff []byte

	//go:embed testdata/before.json
	testBeforeJSON []byte

	//go:embed testdata/after.json
	testAfterJSON []byte
)
