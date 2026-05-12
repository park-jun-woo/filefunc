//ff:func feature=chain type=util control=sequence
//ff:what test: TestExtractCallsForFile
package chain

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func TestExtractCallsForFile(t *testing.T) {
	pf := &model.PythonFile{
		Module: "myapp.service",
		Path:   "myapp/service.py",
		Funcs:  []string{"handle"},
		Calls:  []string{"myapp.utils.validate", "myapp.utils.log_error"},
	}

	calls, err := extractCallsForFile(pf, "", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(calls) != 2 {
		t.Fatalf("calls = %v, want 2 items", calls)
	}
	if calls[0] != "myapp.utils.validate" {
		t.Errorf("calls[0] = %s, want myapp.utils.validate", calls[0])
	}
	if calls[1] != "myapp.utils.log_error" {
		t.Errorf("calls[1] = %s, want myapp.utils.log_error", calls[1])
	}
}
