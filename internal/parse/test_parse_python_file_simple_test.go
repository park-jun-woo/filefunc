//ff:func feature=parse type=util control=sequence
//ff:what test: TestParsePythonFileSimple
package parse

import (
	"testing"
)

func TestParsePythonFileSimple(t *testing.T) {
	pf, err := ParsePythonFile("testdata/simple_func.py")
	if err != nil {
		t.Fatalf("ParsePythonFile failed: %v", err)
	}

	if len(pf.Funcs) != 1 || pf.Funcs[0] != "greet" {
		t.Errorf("Funcs = %v, want [greet]", pf.Funcs)
	}

	if len(pf.Classes) != 0 {
		t.Errorf("Classes = %v, want []", pf.Classes)
	}

	if len(pf.Vars) != 1 || pf.Vars[0] != "MAX_SIZE" {
		t.Errorf("Vars = %v, want [MAX_SIZE]", pf.Vars)
	}

	if pf.Lines != 4 {
		t.Errorf("Lines = %d, want 4", pf.Lines)
	}

	if pf.MaxDepth != 0 {
		t.Errorf("MaxDepth = %d, want 0", pf.MaxDepth)
	}

	if pf.Control != "sequence" {
		t.Errorf("Control = %q, want %q", pf.Control, "sequence")
	}

	if pf.BodyHash == "" {
		t.Error("BodyHash is empty, want non-empty")
	}
}
