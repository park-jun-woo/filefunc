//ff:func feature=validate type=util control=sequence
//ff:what test: TestException_F4_VarWithInit
package validate

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/parse"
)

// F4 exception: var + init() → should not fire
func TestException_F4_VarWithInit(t *testing.T) {
	gf, err := parse.ParseGoFile("testdata/var_with_init.go")
	if err != nil {
		t.Fatal(err)
	}
	expectNoViolation(t, ruleViolations(ExistsWhen, gf, nil, backingF4))
}
