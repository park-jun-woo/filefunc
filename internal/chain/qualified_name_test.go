//ff:func feature=chain type=util control=sequence
//ff:what test: TestQualifiedName
package chain

import "testing"

func TestQualifiedName(t *testing.T) {
	got := qualifiedName("validate", "RuleF1")
	if got != "validate.RuleF1" {
		t.Errorf("qualifiedName = %q, want %q", got, "validate.RuleF1")
	}
}
