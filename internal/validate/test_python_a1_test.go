//ff:func feature=validate type=util control=sequence
//ff:what test: TestPythonA1
package validate

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func TestPythonA1(t *testing.T) {
	pf := mustParsePython(t, "testdata/py_no_annotation.py")
	violations := RunAll([]model.SourceFile{pf}, nil)
	expectViolation(t, violations, "A1")
}
