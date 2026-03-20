//ff:func feature=validate type=util control=sequence
//ff:what test: TestMutest_A1
package validate

import "testing"

// A1
func TestMutest_A1(t *testing.T) {
	expectViolation(t, ruleViolations(ExistsWhen, mustParse(t, "testdata/no_annotation.go"), nil, backingA1f), "A1")
}
