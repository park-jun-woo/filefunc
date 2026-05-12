//ff:func feature=validate type=util control=sequence
//ff:what test: TestPythonClean
package validate

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func TestPythonClean(t *testing.T) {
	pf := mustParsePython(t, "testdata/py_clean.py")
	violations := RunAll([]model.SourceFile{pf}, nil)
	expectNoViolation(t, violations)
}
