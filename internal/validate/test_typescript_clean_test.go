//ff:func feature=validate type=util control=sequence
//ff:what test: TestTypeScriptClean
package validate

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func TestTypeScriptClean(t *testing.T) {
	tf := mustParseTypeScript(t, "testdata/ts_clean.ts")
	violations := RunAll([]model.SourceFile{tf}, nil)
	expectNoViolation(t, violations)
}
