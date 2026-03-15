package validate

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/parse"
)

func TestCheckOneFileOneType_Violation(t *testing.T) {
	gf, err := parse.ParseGoFile("testdata/multi_type.go")
	if err != nil {
		t.Fatal(err)
	}
	violations := CheckOneFileOneType(gf)
	if len(violations) != 1 {
		t.Errorf("expected 1 violation, got %d", len(violations))
	}
	if len(violations) > 0 && violations[0].Rule != "F2" {
		t.Errorf("expected rule F2, got %s", violations[0].Rule)
	}
}
