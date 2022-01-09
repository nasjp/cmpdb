package cmpdb_test

import (
	"testing"

	_ "embed"

	"github.com/google/go-cmp/cmp"
	"github.com/nasjp/cmpdb"
)

func TestSeparateDiff(t *testing.T) {
	t.Parallel()

	got, afterData, err := cmpdb.SeparateDiff(testJSONDiff)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(testBeforeJSON, got); diff != "" {
		t.Errorf("beforeData mismatch (-want +got):\n%s", diff)
	}

	if diff := cmp.Diff(afterData, testAfterJSON); diff != "" {
		t.Errorf("afterData mismatch (-want +got):\n%s", diff)
	}
}
