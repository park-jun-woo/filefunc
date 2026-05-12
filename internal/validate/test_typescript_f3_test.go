//ff:func feature=validate type=util control=iteration dimension=1
//ff:what test: TestTypeScriptF3Exempt — TypeScript F3 면제 확인
package validate

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func TestTypeScriptF3Exempt(t *testing.T) {
	tf := mustParseTypeScript(t, "testdata/ts_multi_method.ts")
	violations := RunAll([]model.SourceFile{tf}, nil)
	for _, v := range violations {
		if v.Rule == "F3" {
			t.Errorf("F3 should be exempt for TypeScript, got violation: %s", v.Message)
		}
	}
}
