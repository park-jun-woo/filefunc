//ff:func feature=validate type=util control=sequence
//ff:what test: TestMutest_A3
package validate

import "testing"

// A3
func TestMutest_A3(t *testing.T) {
	expectViolation(t, ruleViolations(ExistsWhen, mustParse(t, "testdata/no_what.go"), nil, backingA3), "A3")
}
