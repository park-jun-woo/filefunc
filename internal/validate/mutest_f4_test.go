//ff:func feature=validate type=util control=sequence
//ff:what test: TestMutest_F4
package validate

import "testing"

// F4
func TestMutest_F4(t *testing.T) {
	expectViolation(t, ruleViolations(ExistsWhen, mustParse(t, "testdata/init_alone.go"), nil, backingF4), "F4")
	expectNoViolation(t, ruleViolations(ExistsWhen, mustParse(t, "testdata/clean.go"), nil, backingF4))
}
