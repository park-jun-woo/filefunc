//ff:func feature=validate type=util control=sequence
//ff:what test: resolvePythonImport resolves single-dot relative import
package validate

import (
	"path/filepath"
	"testing"
)

func TestResolvePythonImportRelative(t *testing.T) {
	root, err := filepath.Abs("testdata/py_circular")
	if err != nil {
		t.Fatal(err)
	}
	fromFile := filepath.Join(root, "a.py")
	got := resolvePythonImport(fromFile, ".b", root)
	expected := filepath.Join(root, "b.py")
	if got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}
