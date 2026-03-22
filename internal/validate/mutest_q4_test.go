//ff:func feature=validate type=util control=sequence
//ff:what test: TestMutest_Q4
package validate

import "testing"

func TestMutest_Q4(t *testing.T) {
	expectViolation(t, ruleViolations(CheckControlBody, mustParse(t, "testdata/q4_range_long.go"), nil, nil), "Q4")
}
