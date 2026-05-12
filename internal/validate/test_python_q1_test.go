//ff:func feature=validate type=util control=sequence
//ff:what test: TestPythonQ1
package validate

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func TestPythonQ1(t *testing.T) {
	pf := mustParsePython(t, "testdata/py_deep_nesting.py")
	violations := RunAll([]model.SourceFile{pf}, nil)
	expectViolation(t, violations, "Q1")
}
