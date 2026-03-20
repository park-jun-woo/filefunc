//ff:func feature=report type=util control=sequence
//ff:what test: TestFormatJSONEmpty
package report

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
)

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
