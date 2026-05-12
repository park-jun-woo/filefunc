//ff:func feature=parse type=util control=iteration dimension=1
//ff:what test: TestParsePythonFileClass
package parse

import (
	"testing"
)

func TestParsePythonFileClass(t *testing.T) {
	pf, err := ParsePythonFile("testdata/with_class.py")
	if err != nil {
		t.Fatalf("ParsePythonFile failed: %v", err)
	}

	if len(pf.Classes) != 1 || pf.Classes[0] != "Server" {
		t.Errorf("Classes = %v, want [Server]", pf.Classes)
	}

	if len(pf.Methods) != 2 {
		t.Errorf("Methods count = %d, want 2", len(pf.Methods))
	}

	if !pf.HasInitMethod {
		t.Error("HasInitMethod = false, want true")
	}

	if len(pf.Funcs) != 0 {
		t.Errorf("Funcs = %v, want []", pf.Funcs)
	}

	wantMethods := map[string]bool{"Server.start": true, "Server.stop": true}
	for _, m := range pf.Methods {
		if !wantMethods[m] {
			t.Errorf("unexpected method %q", m)
		}
	}
}
