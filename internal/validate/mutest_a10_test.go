//ff:func feature=validate type=util control=sequence
//ff:what test: TestMutest_A10
package validate

import "testing"

// A10
func TestMutest_A10(t *testing.T) {
	expectViolation(t, ruleViolations(ControlMatch, mustParse(t, "testdata/selection_no_switch.go"), nil, backingA10), "A10")
}
