//ff:func feature=validate type=util control=sequence
//ff:what test: TestPythonA6
package validate

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func TestPythonA6(t *testing.T) {
	pf := mustParsePython(t, "testdata/py_annotation_after_import.py")
	violations := RunAll([]model.SourceFile{pf}, nil)
	expectViolation(t, violations, "A6")
}
