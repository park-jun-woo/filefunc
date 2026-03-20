//ff:func feature=validate type=util control=sequence
//ff:what test: TestType_A1_Present
package validate

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/parse"
)

// A1: type-only file with //ff:type → should pass
func TestType_A1_Present(t *testing.T) {
	gf, err := parse.ParseGoFile("testdata/type_with_annotation.go")
	if err != nil {
		t.Fatal(err)
	}
	expectNoViolation(t, ruleViolations(ExistsWhen, gf, nil, backingA1t))
}
