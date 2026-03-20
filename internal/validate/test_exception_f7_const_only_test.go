//ff:func feature=validate type=util control=sequence
//ff:what test: TestException_F7_ConstOnly
package validate

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
	"github.com/park-jun-woo/filefunc/internal/parse"
)

// F7: const-only file → F1 should not fire (defeated by IsConstOnlyDefeater)
func TestException_F7_ConstOnly(t *testing.T) {
	gf, err := parse.ParseGoFile("testdata/const_only.go")
	if err != nil {
		t.Fatal(err)
	}
	violations := RunAll([]*model.GoFile{gf}, nil)
	expectNoViolation(t, violations)
}
