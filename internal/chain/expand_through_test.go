//ff:func feature=chain type=util control=sequence
//ff:what test: TestExpandThrough
package chain

import "testing"

func TestExpandThrough(t *testing.T) {
	g := buildTestGraph()
	// From Caller's children, expand their children
	result := ExpandThrough(g.Children["testdata.Caller"], func(c string) []string {
		return g.Children[c]
	})
	if len(result) != 1 || result[0] != "testdata.Leaf" {
		t.Errorf("ExpandThrough = %v, want [testdata.Leaf]", result)
	}
}
