//ff:func feature=chain type=util control=iteration dimension=1
//ff:what test: TestTraverseChon_Chon1
package chain

import "testing"

func TestTraverseChon_Chon1(t *testing.T) {
	g := buildTestGraph()
	results := TraverseChon(g, "testdata.HelperA", 1)

	hasChild := false
	hasParent := false
	for _, r := range results {
		if r.Name == "testdata.Leaf" && r.Chon == 1 && r.Rel == "calls" {
			hasChild = true
		}
		if r.Name == "testdata.Caller" && r.Chon == 1 && r.Rel == "called-by" {
			hasParent = true
		}
	}
	if !hasChild {
		t.Error("missing chon=1 child: testdata.Leaf")
	}
	if !hasParent {
		t.Error("missing chon=1 parent: testdata.Caller")
	}
}
