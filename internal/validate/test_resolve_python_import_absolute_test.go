//ff:func feature=validate type=util control=sequence
//ff:what test: resolvePythonImport returns empty for absolute import
package validate

import (
	"path/filepath"
	"testing"
)

func TestResolvePythonImportAbsolute(t *testing.T) {
	root, err := filepath.Abs("testdata/py_circular")
	if err != nil {
		t.Fatal(err)
	}
	fromFile := filepath.Join(root, "a.py")
	got := resolvePythonImport(fromFile, "os", root)
	if got != "" {
		t.Errorf("expected empty string for absolute import, got %s", got)
	}
}
