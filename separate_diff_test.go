package cmpdb_test

import (
	"testing"

	_ "embed"

	"github.com/google/go-cmp/cmp"
	"github.com/nasjp/cmpdb"
)

func TestSeparateDiff(t *testing.T) {
	t.Parallel()

	gotBeforeData, gotAfterData, err := cmpdb.SeparateDiff(testJSONDiff)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(testBeforeJSON, gotBeforeData); diff != "" {
		t.Errorf("beforeData mismatch (-want +got):\n%s", diff)
	}

	if diff := cmp.Diff(gotAfterData, testAfterJSON); diff != "" {
		t.Errorf("afterData mismatch (-want +got):\n%s", diff)
	}
}
