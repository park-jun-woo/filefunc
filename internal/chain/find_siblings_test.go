//ff:func feature=chain type=util control=sequence
//ff:what test: TestFindSiblings
package chain

import "testing"

func TestFindSiblings(t *testing.T) {
	g := buildTestGraph()
	siblings := FindSiblings(g, "testdata.HelperA")
	if len(siblings) != 1 || siblings[0] != "testdata.HelperB" {
		t.Errorf("FindSiblings = %v, want [testdata.HelperB]", siblings)
	}
}
