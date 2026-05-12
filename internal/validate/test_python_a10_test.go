//ff:func feature=validate type=util control=sequence
//ff:what test: TestPythonA10
package validate

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func TestPythonA10(t *testing.T) {
	pf := mustParsePython(t, "testdata/py_selection_no_match.py")
	violations := RunAll([]model.SourceFile{pf}, nil)
	expectViolation(t, violations, "A10")
}
