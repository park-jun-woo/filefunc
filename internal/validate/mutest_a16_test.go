//ff:func feature=validate type=util control=sequence
//ff:what test: TestMutest_A16
package validate

import "testing"

// A16
func TestMutest_A16(t *testing.T) {
	expectViolation(t, ruleViolations(ValidDimension, mustParse(t, "testdata/bad_dimension_value.go"), nil, nil), "A16")
	expectNoViolation(t, ruleViolations(ValidDimension, mustParse(t, "testdata/clean.go"), nil, nil))
}
