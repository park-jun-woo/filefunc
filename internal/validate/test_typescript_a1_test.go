//ff:func feature=validate type=util control=sequence
//ff:what test: TestTypeScriptA1 — func 파일에 //ff:func 없으면 A1 위반
package validate

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func TestTypeScriptA1(t *testing.T) {
	tf := mustParseTypeScript(t, "testdata/ts_a1_missing.ts")
	violations := RunAll([]model.SourceFile{tf}, nil)
	expectViolation(t, violations, "A1")
}
