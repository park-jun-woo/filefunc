//ff:func feature=validate type=util control=sequence
//ff:what test: TestException_F6_ParamType
package validate

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/parse"
)

// F6: func + unexported param type → F2 should not fire
func TestException_F6_ParamType(t *testing.T) {
	gf, err := parse.ParseGoFile("testdata/func_with_param_type.go")
	if err != nil {
		t.Fatal(err)
	}
	expectNoViolation(t, ruleViolations(CountMax, gf, nil, backingF2))
}
