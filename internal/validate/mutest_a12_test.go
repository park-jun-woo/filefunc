//ff:func feature=validate type=util control=sequence
//ff:what test: TestMutest_A12
package validate

import "testing"

// A12
func TestMutest_A12(t *testing.T) {
	expectViolation(t, ruleViolations(ControlMatch, mustParse(t, "testdata/sequence_with_loop.go"), nil, backingA12), "A12")
}
