//ff:func feature=validate type=util control=sequence
//ff:what test: TestType_A1_Missing
package validate

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/parse"
)

// A1: type-only file without //ff:type → should fire
func TestType_A1_Missing(t *testing.T) {
	gf, err := parse.ParseGoFile("testdata/type_no_annotation.go")
	if err != nil {
		t.Fatal(err)
	}
	expectViolation(t, ruleViolations(ExistsWhen, gf, nil, backingA1t), "A1")
}
