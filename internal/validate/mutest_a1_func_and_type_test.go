//ff:func feature=validate type=util control=sequence
//ff:what test: TestMutest_A1_FuncAndType
package validate

import "testing"

// A1: file with func+type should check both annotations
func TestMutest_A1_FuncAndType(t *testing.T) {
	gf := mustParse(t, "testdata/sample_with_func_and_type.go")
	// has //ff:func but not //ff:type — should fire A1 for type
	violations := ruleViolations(ExistsWhen, gf, nil, backingA1t)
	expectViolation(t, violations, "A1")
}
