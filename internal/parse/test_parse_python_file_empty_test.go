//ff:func feature=parse type=util control=sequence
//ff:what test: TestParsePythonFileEmpty
package parse

import (
	"testing"
)

func TestParsePythonFileEmpty(t *testing.T) {
	pf, err := ParsePythonFile("testdata/empty.py")
	if err != nil {
		t.Fatalf("ParsePythonFile failed: %v", err)
	}

	if len(pf.Funcs) != 0 {
		t.Errorf("Funcs = %v, want []", pf.Funcs)
	}

	if len(pf.Classes) != 0 {
		t.Errorf("Classes = %v, want []", pf.Classes)
	}

	if pf.MaxDepth != 0 {
		t.Errorf("MaxDepth = %d, want 0", pf.MaxDepth)
	}

	if pf.Lines != 0 {
		t.Errorf("Lines = %d, want 0", pf.Lines)
	}

	if pf.BodyHash != "" {
		t.Errorf("BodyHash = %q, want empty", pf.BodyHash)
	}
}
