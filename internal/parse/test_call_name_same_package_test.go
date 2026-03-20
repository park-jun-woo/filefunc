//ff:func feature=parse type=util control=iteration dimension=1
//ff:what test: TestCallName_SamePackage
package parse

import (
	"testing"
)

func TestCallName_SamePackage(t *testing.T) {
	projFuncs := map[string]string{"HelperA": "testdata", "Leaf": "testdata"}
	projImports := map[string]string{}

	calls, err := ExtractCalls("testdata/caller.go", "github.com/nonexistent", projFuncs, "testdata")
	if err != nil {
		t.Fatalf("ExtractCalls failed: %v", err)
	}

	found := make(map[string]bool)
	for _, c := range calls {
		found[c] = true
	}
	if !found["testdata.HelperA"] {
		t.Errorf("missing call testdata.HelperA, got %v", calls)
	}
	_ = projImports
}
