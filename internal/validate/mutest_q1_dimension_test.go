//ff:func feature=validate type=util control=sequence
//ff:what test: TestMutest_Q1_Dimension
package validate

import "testing"

// Q1 dimension
func TestMutest_Q1_Dimension(t *testing.T) {
	expectNoViolation(t, ruleViolations(CheckDepthLimit, mustParse(t, "testdata/dimension2_depth3.go"), nil, nil))
}
