//ff:func feature=validate type=util control=sequence
//ff:what test: TestMutest_Q4_NestedOK
package validate

import "testing"

func TestMutest_Q4_NestedOK(t *testing.T) {
	violations := ruleViolations(CheckControlBody, mustParse(t, "testdata/q4_range_nested_ok.go"), nil, nil)
	if len(violations) > 0 {
		t.Errorf("expected no Q4 violation for nested range with low pure lines, got %v", violations)
	}
}
