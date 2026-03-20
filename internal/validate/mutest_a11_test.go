//ff:func feature=validate type=util control=sequence
//ff:what test: TestMutest_A11
package validate

import "testing"

// A11
func TestMutest_A11(t *testing.T) {
	expectViolation(t, ruleViolations(ControlMatch, mustParse(t, "testdata/iteration_no_loop.go"), nil, backingA11), "A11")
}
