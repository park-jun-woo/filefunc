//ff:func feature=report type=util control=sequence
//ff:what test: TestFormatJSONWithViolations
package report

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
)

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
