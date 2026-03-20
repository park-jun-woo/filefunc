//ff:func feature=validate type=util control=sequence
//ff:what test: TestMutest_A15
package validate

import "testing"

// A15
func TestMutest_A15(t *testing.T) {
	expectViolation(t, ruleViolations(ExistsWhen, mustParse(t, "testdata/iter_no_dimension.go"), nil, backingA15), "A15")
	expectNoViolation(t, ruleViolations(ExistsWhen, mustParse(t, "testdata/clean.go"), nil, backingA15))
}
