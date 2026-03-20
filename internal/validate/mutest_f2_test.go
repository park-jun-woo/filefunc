//ff:func feature=validate type=util control=sequence
//ff:what test: TestMutest_F2
package validate

import "testing"

// F2
func TestMutest_F2(t *testing.T) {
	expectViolation(t, ruleViolations(CountMax, mustParse(t, "testdata/multi_type.go"), nil, backingF2), "F2")
}
