//ff:func feature=validate type=util control=sequence
//ff:what test: TestTypeScriptA6 — 어노테이션이 import 뒤에 위치하면 A6 위반
package validate

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func TestTypeScriptA6(t *testing.T) {
	tf := mustParseTypeScript(t, "testdata/ts_a6_bad.ts")
	violations := RunAll([]model.SourceFile{tf}, nil)
	expectViolation(t, violations, "A6")
}
