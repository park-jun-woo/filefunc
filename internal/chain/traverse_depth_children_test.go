//ff:func feature=chain type=util control=iteration dimension=1
//ff:what test: TestTraverseDepth_Children
package chain

import "testing"

func TestTraverseDepth_Children(t *testing.T) {
	g := buildTestGraph()
	results := TraverseDepth(g, "testdata.Caller", "calls", 3)

	names := make(map[string]bool)
	for _, r := range results {
		names[r.Name] = true
	}
	if !names["testdata.HelperA"] || !names["testdata.HelperB"] || !names["testdata.Leaf"] {
		t.Errorf("TraverseDepth children = %v, want HelperA, HelperB, Leaf", results)
	}
}
