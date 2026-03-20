//ff:func feature=validate type=util control=sequence
//ff:what test: TestMutest_F3
package validate

import "testing"

// F3
func TestMutest_F3(t *testing.T) {
	expectViolation(t, ruleViolations(CountMax, mustParse(t, "testdata/multi_method.go"), nil, backingF3), "F3")
}
