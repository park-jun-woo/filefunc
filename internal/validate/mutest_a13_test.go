//ff:func feature=validate type=util control=sequence
//ff:what test: TestMutest_A13
package validate

import "testing"

// A13
func TestMutest_A13(t *testing.T) {
	expectViolation(t, ruleViolations(ControlMatch, mustParse(t, "testdata/control_selection_with_loop.go"), nil, backingA13), "A13")
	expectNoViolation(t, ruleViolations(ControlMatch, mustParse(t, "testdata/clean.go"), nil, backingA13))
}
