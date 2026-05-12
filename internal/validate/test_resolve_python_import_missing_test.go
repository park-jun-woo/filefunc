//ff:func feature=validate type=util control=sequence
//ff:what test: resolvePythonImport returns empty for nonexistent file
package validate

import (
	"path/filepath"
	"testing"
)

func TestResolvePythonImportMissing(t *testing.T) {
	root, err := filepath.Abs("testdata/py_circular")
	if err != nil {
		t.Fatal(err)
	}
	fromFile := filepath.Join(root, "a.py")
	got := resolvePythonImport(fromFile, ".nonexistent", root)
	if got != "" {
		t.Errorf("expected empty string for missing file, got %s", got)
	}
}
