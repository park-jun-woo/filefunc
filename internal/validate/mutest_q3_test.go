//ff:func feature=validate type=util control=sequence
//ff:what test: TestMutest_Q3
package validate

import "testing"

// Q3
func TestMutest_Q3(t *testing.T) {
	expectViolation(t, ruleViolations(CheckFuncLines, mustParse(t, "testdata/medium_func.go"), nil, nil), "Q3")
}
