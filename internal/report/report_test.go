package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func TestFormatTextNoViolations(t *testing.T) {
	var buf bytes.Buffer
	FormatText(&buf, nil)
	if !strings.Contains(buf.String(), "No violations found.") {
		t.Errorf("got %q", buf.String())
	}
}

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

func TestFormatJSONEmpty(t *testing.T) {
	var buf bytes.Buffer
	if err := FormatJSON(&buf, nil); err != nil {
		t.Fatal(err)
	}
	var result []model.Violation
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty array, got %d", len(result))
	}
}

func TestFormatJSONWithViolations(t *testing.T) {
	var buf bytes.Buffer
	vs := []model.Violation{{File: "a.go", Rule: "F1", Level: "ERROR", Message: "msg"}}
	if err := FormatJSON(&buf, vs); err != nil {
		t.Fatal(err)
	}
	var result []model.Violation
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 || result[0].Rule != "F1" {
		t.Errorf("got %+v", result)
	}
}
