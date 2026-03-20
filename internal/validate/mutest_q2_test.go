//ff:func feature=validate type=util control=sequence
//ff:what test: TestMutest_Q2
package validate

import "testing"

// Q2
func TestMutest_Q2(t *testing.T) {
	expectViolation(t, ruleViolations(CheckFuncLines, mustParse(t, "testdata/long_func.go"), nil, nil), "Q2")
}
