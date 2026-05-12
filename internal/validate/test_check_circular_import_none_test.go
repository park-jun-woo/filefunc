//ff:func feature=validate type=util control=sequence
//ff:what test: CheckCircularImport returns no violation for acyclic imports
package validate

import (
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/filefunc/internal/parse"
)

func TestCheckCircularImportNone(t *testing.T) {
	root, err := filepath.Abs("testdata/py_no_circular")
	if err != nil {
		t.Fatal(err)
	}
	paths := []string{
		filepath.Join(root, "a.py"),
		filepath.Join(root, "b.py"),
		filepath.Join(root, "c.py"),
	}
	pyFiles, err := parse.ParsePythonFiles(paths)
	if err != nil {
		t.Fatal(err)
	}
	violations := CheckCircularImport(pyFiles, root)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %v", violations)
	}
}
