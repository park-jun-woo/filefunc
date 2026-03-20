//ff:func feature=validate type=util control=sequence
//ff:what test: TestClean_AllRules
package validate

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
	"github.com/park-jun-woo/filefunc/internal/parse"
)

// --- Clean: all rules pass ---

func TestClean_AllRules(t *testing.T) {
	gf, err := parse.ParseGoFile("testdata/clean.go")
	if err != nil {
		t.Fatal(err)
	}
	violations := RunAll([]*model.GoFile{gf}, nil)
	expectNoViolation(t, violations)
}
