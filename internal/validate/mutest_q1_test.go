//ff:func feature=validate type=util control=sequence
//ff:what test: TestMutest_Q1
package validate

import "testing"

// Q1
func TestMutest_Q1(t *testing.T) {
	expectViolation(t, ruleViolations(CheckDepthLimit, mustParse(t, "testdata/deep_nesting.go"), nil, nil), "Q1")
	expectNoViolation(t, ruleViolations(CheckDepthLimit, mustParse(t, "testdata/clean.go"), nil, nil))
}
