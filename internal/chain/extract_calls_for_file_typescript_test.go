//ff:func feature=chain type=util control=sequence
//ff:what test: TestExtractCallsForFileTypeScript
package chain

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func TestExtractCallsForFileTypeScript(t *testing.T) {
	tf := &model.TypeScriptFile{
		Module: "src/service",
		Path:   "src/service.ts",
		Funcs:  []string{"handleRequest"},
		Calls:  []string{"src/utils.validate", "src/utils.logError"},
	}

	calls, err := extractCallsForFile(tf, "", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(calls) != 2 {
		t.Fatalf("calls = %v, want 2 items", calls)
	}
	if calls[0] != "src/utils.validate" {
		t.Errorf("calls[0] = %s, want src/utils.validate", calls[0])
	}
	if calls[1] != "src/utils.logError" {
		t.Errorf("calls[1] = %s, want src/utils.logError", calls[1])
	}
}
