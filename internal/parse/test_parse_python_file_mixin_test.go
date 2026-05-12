//ff:func feature=parse type=util control=iteration dimension=1
//ff:what test: TestParsePythonFileMixin
package parse

import (
	"testing"
)

func TestParsePythonFileMixin(t *testing.T) {
	pf, err := ParsePythonFile("testdata/mixin.py")
	if err != nil {
		t.Fatalf("ParsePythonFile failed: %v", err)
	}

	if len(pf.Classes) != 2 {
		t.Errorf("Classes count = %d, want 2", len(pf.Classes))
	}

	wantClasses := map[string]bool{"LogMixin": true, "AuthMixin": true}
	for _, c := range pf.Classes {
		if !wantClasses[c] {
			t.Errorf("unexpected class %q", c)
		}
	}

	if len(pf.Methods) != 2 {
		t.Errorf("Methods count = %d, want 2", len(pf.Methods))
	}

	wantMethods := map[string]bool{"LogMixin.log": true, "AuthMixin.authenticate": true}
	for _, m := range pf.Methods {
		if !wantMethods[m] {
			t.Errorf("unexpected method %q", m)
		}
	}

	if pf.HasInitMethod {
		t.Error("HasInitMethod = true, want false")
	}
}
