//ff:func feature=chain type=util control=iteration dimension=1
//ff:what test: TestTraverseChon_Chon2
package chain

import "testing"

func TestTraverseChon_Chon2(t *testing.T) {
	g := buildTestGraph()
	results := TraverseChon(g, "testdata.HelperA", 2)

	hasSibling := false
	for _, r := range results {
		if r.Name == "testdata.HelperB" && r.Chon == 2 && r.Rel == "co-called" {
			hasSibling = true
		}
	}
	if !hasSibling {
		t.Error("missing chon=2 sibling: testdata.HelperB")
	}
}
