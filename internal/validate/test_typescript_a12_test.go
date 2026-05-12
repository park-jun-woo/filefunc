//ff:func feature=validate type=util control=sequence
//ff:what test: TestTypeScriptA12 — control=sequence인데 for 존재하면 A12 위반
package validate

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func TestTypeScriptA12(t *testing.T) {
	tf := mustParseTypeScript(t, "testdata/ts_a12_bad.ts")
	violations := RunAll([]model.SourceFile{tf}, nil)
	expectViolation(t, violations, "A12")
}
