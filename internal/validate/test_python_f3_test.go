//ff:func feature=validate type=util control=iteration dimension=1
//ff:what test: TestPythonF3Exempt — Python F3 면제 확인
package validate

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func TestPythonF3Exempt(t *testing.T) {
	pf := mustParsePython(t, "testdata/py_multi_method.py")
	violations := RunAll([]model.SourceFile{pf}, nil)
	for _, v := range violations {
		if v.Rule == "F3" {
			t.Errorf("F3 should be exempt for Python, got violation: %s", v.Message)
		}
	}
}
