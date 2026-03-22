//ff:func feature=validate type=util control=sequence
//ff:what test: TestMutest_Q3
package validate

import (
	"testing"
)

// Q3: sequence func > 100 lines → ERROR
func TestMutest_Q3(t *testing.T) {
	violations := ruleViolations(CheckFuncLines, mustParse(t, "testdata/medium_func.go"), nil, nil)
	expectViolation(t, violations, "Q3")
	if len(violations) > 0 && violations[0].Level != "ERROR" {
		t.Errorf("expected Q3 level ERROR, got %q", violations[0].Level)
	}
}
