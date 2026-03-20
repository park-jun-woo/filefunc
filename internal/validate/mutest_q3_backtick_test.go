//ff:func feature=validate type=util control=sequence
//ff:what test: TestMutest_Q3_Backtick
package validate

import (
	"strings"
	"testing"
)

// Q3 backtick hint
func TestMutest_Q3_Backtick(t *testing.T) {
	violations := ruleViolations(CheckFuncLines, mustParse(t, "testdata/q3_backtick.go"), nil, nil)
	expectViolation(t, violations, "Q3")
	if len(violations) > 0 && !strings.Contains(violations[0].Message, "var-only file") {
		t.Errorf("expected backtick hint in message, got %q", violations[0].Message)
	}
}
