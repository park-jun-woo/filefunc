//ff:func feature=chain type=util control=iteration dimension=1
//ff:what test: TestTraverseDepth_Parents
package chain

import "testing"

func TestTraverseDepth_Parents(t *testing.T) {
	g := buildTestGraph()
	results := TraverseDepth(g, "testdata.Leaf", "called-by", 3)

	names := make(map[string]bool)
	for _, r := range results {
		names[r.Name] = true
	}
	if !names["testdata.HelperA"] || !names["testdata.Caller"] {
		t.Errorf("TraverseDepth parents = %v, want HelperA, Caller", results)
	}
}
