package validate

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/parse"
)

func TestCheckFuncLines_Clean(t *testing.T) {
	gf, err := parse.ParseGoFile("testdata/clean.go")
	if err != nil {
		t.Fatal(err)
	}
	violations := CheckFuncLines(gf)
	if len(violations) != 0 {
		t.Errorf("expected 0 violations, got %d", len(violations))
	}
}
