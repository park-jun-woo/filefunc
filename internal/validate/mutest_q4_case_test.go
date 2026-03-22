//ff:func feature=validate type=util control=sequence
//ff:what test: TestMutest_Q4_Case
package validate

import "testing"

func TestMutest_Q4_Case(t *testing.T) {
	expectViolation(t, ruleViolations(CheckControlBody, mustParse(t, "testdata/q4_case_long.go"), nil, nil), "Q4")
}
