//ff:func feature=validate type=util control=sequence
//ff:what test: TestMutest_F1
package validate

import "testing"

// F1
func TestMutest_F1(t *testing.T) {
	expectViolation(t, ruleViolations(CountMax, mustParse(t, "testdata/multi_func.go"), nil, backingF1), "F1")
	expectNoViolation(t, ruleViolations(CountMax, mustParse(t, "testdata/clean.go"), nil, backingF1))
}
