//ff:func feature=validate type=util control=sequence
//ff:what test: TestMutest_A9
package validate

import "testing"

// A9
func TestMutest_A9(t *testing.T) {
	expectViolation(t, ruleViolations(ExistsWhen, mustParse(t, "testdata/no_control.go"), nil, backingA9), "A9")
}
