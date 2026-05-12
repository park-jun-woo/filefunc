//ff:func feature=validate type=util control=iteration dimension=1
//ff:what test: CheckCircularImport detects I1 violation for circular imports
package validate

import (
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/filefunc/internal/parse"
)

func TestCheckCircularImportFound(t *testing.T) {
	root, err := filepath.Abs("testdata/py_circular")
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
	if len(violations) == 0 {
		t.Error("expected I1 violation for circular import, got none")
		return
	}
	for _, v := range violations {
		if v.Rule != "I1" {
			t.Errorf("expected rule I1, got %s", v.Rule)
		}
	}
}
