//ff:func feature=chain type=util control=sequence
//ff:what test: TestFilterByPackage
package chain

import "testing"

func TestFilterByPackage(t *testing.T) {
	results := []ChonResult{
		{"validate.RuleF1", 1, "child"},
		{"chain.Build", 2, "grandchild"},
	}
	got := FilterByPackage(results, "validate")
	if len(got) != 1 || got[0].Name != "validate.RuleF1" {
		t.Errorf("FilterByPackage = %v", got)
	}
}
