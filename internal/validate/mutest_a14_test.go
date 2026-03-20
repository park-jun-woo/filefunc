//ff:func feature=validate type=util control=sequence
//ff:what test: TestMutest_A14
package validate

import "testing"

// A14
func TestMutest_A14(t *testing.T) {
	expectViolation(t, ruleViolations(ControlMatch, mustParse(t, "testdata/control_iteration_with_switch.go"), nil, backingA14), "A14")
	expectNoViolation(t, ruleViolations(ControlMatch, mustParse(t, "testdata/clean.go"), nil, backingA14))
}
