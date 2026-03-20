//ff:func feature=parse type=util control=sequence
//ff:what test: TestCallName_ExternalIgnored
package parse

import (
	"testing"
)

func TestCallName_ExternalIgnored(t *testing.T) {
	projFuncs := map[string]string{}
	calls, err := ExtractCalls("testdata/caller.go", "github.com/nonexistent", projFuncs, "testdata")
	if err != nil {
		t.Fatalf("ExtractCalls failed: %v", err)
	}
	if len(calls) != 0 {
		t.Errorf("expected no calls for unknown funcs, got %v", calls)
	}
}
