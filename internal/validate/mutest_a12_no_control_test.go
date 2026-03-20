//ff:func feature=validate type=util control=sequence
//ff:what test: TestMutest_A12_NoControl
package validate

import "testing"

// A12: control="" should NOT fire A12 (A9 handles it)
func TestMutest_A12_NoControl(t *testing.T) {
	expectNoViolation(t, ruleViolations(ControlMatch, mustParse(t, "testdata/no_control.go"), nil, backingA12))
}
