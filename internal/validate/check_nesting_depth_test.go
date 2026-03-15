package validate

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/parse"
)

func TestCheckNestingDepth_Violation(t *testing.T) {
	gf, err := parse.ParseGoFile("testdata/deep_nesting.go")
	if err != nil {
		t.Fatal(err)
	}
	violations := CheckNestingDepth(gf)
	if len(violations) != 1 {
		t.Errorf("expected 1 violation, got %d", len(violations))
	}
	if len(violations) > 0 && violations[0].Rule != "Q1" {
		t.Errorf("expected rule Q1, got %s", violations[0].Rule)
	}
}

func TestCheckNestingDepth_Clean(t *testing.T) {
	gf, err := parse.ParseGoFile("testdata/clean.go")
	if err != nil {
		t.Fatal(err)
	}
	violations := CheckNestingDepth(gf)
	if len(violations) != 0 {
		t.Errorf("expected 0 violations, got %d", len(violations))
	}
}
