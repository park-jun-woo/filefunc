//ff:func feature=validate type=util control=sequence
//ff:what test: TestTypeScriptF3
package validate

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func TestTypeScriptF3(t *testing.T) {
	tf := mustParseTypeScript(t, "testdata/ts_multi_method.ts")
	violations := RunAll([]model.SourceFile{tf}, nil)
	expectViolation(t, violations, "F3")
}
