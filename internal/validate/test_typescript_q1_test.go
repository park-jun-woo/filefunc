//ff:func feature=validate type=util control=sequence
//ff:what test: TestTypeScriptQ1
package validate

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func TestTypeScriptQ1(t *testing.T) {
	tf := mustParseTypeScript(t, "testdata/ts_deep_nesting.ts")
	violations := RunAll([]model.SourceFile{tf}, nil)
	expectViolation(t, violations, "Q1")
}
