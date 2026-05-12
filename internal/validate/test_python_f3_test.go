//ff:func feature=validate type=util control=sequence
//ff:what test: TestPythonF3
package validate

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func TestPythonF3(t *testing.T) {
	pf := mustParsePython(t, "testdata/py_multi_method.py")
	violations := RunAll([]model.SourceFile{pf}, nil)
	expectViolation(t, violations, "F3")
}
