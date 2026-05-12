//ff:func feature=validate type=util control=sequence
//ff:what test: TestTypeScriptA10 — control=selection인데 switch 없으면 A10 위반
package validate

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func TestTypeScriptA10(t *testing.T) {
	tf := mustParseTypeScript(t, "testdata/ts_a10_bad.ts")
	violations := RunAll([]model.SourceFile{tf}, nil)
	expectViolation(t, violations, "A10")
}
