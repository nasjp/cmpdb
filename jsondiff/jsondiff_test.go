package jsondiff_test

import (
	"testing"

	_ "embed"

	"github.com/nasjp/cmpdb/jsondiff"
)

//go:embed testdata/diff.jsondiff
var testJSONDiff []byte

func TestUnmarshal(t *testing.T) {
	t.Parallel()

	var beforeDist, afterDist interface{}

	if err := jsondiff.Unmarshal(testJSONDiff, &beforeDist, &afterDist); err != nil {
		t.Fatal(err)
	}

	t.Logf("%#v", beforeDist)
	t.Logf("%#v", afterDist)
}
