//ff:func feature=validate type=util control=sequence
//ff:what test: CheckCircularImport ignores external imports in graph
package validate

import (
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func TestCheckCircularImportExternal(t *testing.T) {
	root, err := filepath.Abs("testdata/py_no_circular")
	if err != nil {
		t.Fatal(err)
	}
	files := []*model.PythonFile{
		{
			Path:          filepath.Join(root, "b.py"),
			ModuleImports: []string{"os", "sys", "json"},
		},
	}
	violations := CheckCircularImport(files, root)
	if len(violations) != 0 {
		t.Errorf("expected no violations for external imports, got %v", violations)
	}
}
