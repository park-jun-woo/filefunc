//ff:func feature=report type=util control=sequence
//ff:what test: TestFormatTextNoViolations
package report

import (
	"bytes"
	"strings"
	"testing"
)

func TestFormatTextNoViolations(t *testing.T) {
	var buf bytes.Buffer
	FormatText(&buf, nil)
	if !strings.Contains(buf.String(), "No violations found.") {
		t.Errorf("got %q", buf.String())
	}
}
