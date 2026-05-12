//ff:func feature=parse type=util control=sequence
//ff:what test: TestParsePythonFileNested
package parse

import (
	"testing"
)

func TestParsePythonFileNested(t *testing.T) {
	pf, err := ParsePythonFile("testdata/nested_control.py")
	if err != nil {
		t.Fatalf("ParsePythonFile failed: %v", err)
	}

	if pf.MaxDepth != 3 {
		t.Errorf("MaxDepth = %d, want 3", pf.MaxDepth)
	}

	if pf.Control != "iteration" {
		t.Errorf("Control = %q, want %q", pf.Control, "iteration")
	}

	if !pf.HasLoopAtDepth1 {
		t.Error("HasLoopAtDepth1 = false, want true")
	}

	if pf.HasMatchAtDepth1 {
		t.Error("HasMatchAtDepth1 = true, want false")
	}

	if len(pf.Funcs) != 1 || pf.Funcs[0] != "process" {
		t.Errorf("Funcs = %v, want [process]", pf.Funcs)
	}
}
