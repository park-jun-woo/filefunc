//ff:func feature=validate type=util control=sequence
//ff:what test: TestPythonA12
package validate

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func TestPythonA12(t *testing.T) {
	pf := mustParsePython(t, "testdata/py_sequence_with_for.py")
	violations := RunAll([]model.SourceFile{pf}, nil)
	expectViolation(t, violations, "A12")
}
