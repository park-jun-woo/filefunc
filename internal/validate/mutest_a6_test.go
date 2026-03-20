//ff:func feature=validate type=util control=sequence
//ff:what test: TestMutest_A6
package validate

import "testing"

// A6
func TestMutest_A6(t *testing.T) {
	expectViolation(t, ruleViolations(AnnotationAtTop, mustParse(t, "testdata/annotation_after_func.go"), nil, nil), "A6")
}
