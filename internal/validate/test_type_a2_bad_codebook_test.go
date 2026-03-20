//ff:func feature=validate type=util control=sequence
//ff:what test: TestType_A2_BadCodebook
package validate

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
	"github.com/park-jun-woo/filefunc/internal/parse"
)

// A2: //ff:type with bad codebook value → should fire
func TestType_A2_BadCodebook(t *testing.T) {
	gf, err := parse.ParseGoFile("testdata/type_bad_codebook.go")
	if err != nil {
		t.Fatal(err)
	}
	cb := &model.Codebook{
		Required: map[string]map[string]string{
			"feature": {"validate": "", "parse": ""},
			"type":    {"rule": "", "model": ""},
		},
	}
	expectViolation(t, ruleViolations(InCodebook, gf, cb, backingA2), "A2")
}
