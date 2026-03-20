//ff:func feature=report type=util control=sequence
//ff:what test: TestFormatTextWithViolations
package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func TestFormatTextWithViolations(t *testing.T) {
	var buf bytes.Buffer
	vs := []model.Violation{
		{File: "a.go", Rule: "F1", Level: "ERROR", Message: "too many funcs"},
		{File: "b.go", Rule: "Q1", Level: "WARNING", Message: "deep nesting"},
	}
	FormatText(&buf, vs)
	out := buf.String()
	if !strings.Contains(out, "[ERROR] F1") {
		t.Error("missing ERROR F1")
	}
	if !strings.Contains(out, "[WARNING] Q1") {
		t.Error("missing WARNING Q1")
	}
	if !strings.Contains(out, "2 violation(s) found.") {
		t.Error("missing count")
	}
}
