package validate

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/parse"
)

func TestCheckInitStandalone_Violation(t *testing.T) {
	gf, err := parse.ParseGoFile("testdata/init_alone.go")
	if err != nil {
		t.Fatal(err)
	}
	violations := CheckInitStandalone(gf)
	if len(violations) != 1 {
		t.Errorf("expected 1 violation, got %d", len(violations))
	}
	if len(violations) > 0 && violations[0].Rule != "F4" {
		t.Errorf("expected rule F4, got %s", violations[0].Rule)
	}
}

func TestCheckInitStandalone_Clean(t *testing.T) {
	gf, err := parse.ParseGoFile("testdata/clean.go")
	if err != nil {
		t.Fatal(err)
	}
	violations := CheckInitStandalone(gf)
	if len(violations) != 0 {
		t.Errorf("expected 0 violations, got %d", len(violations))
	}
}
