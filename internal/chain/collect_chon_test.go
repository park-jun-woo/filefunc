//ff:func feature=chain type=util control=sequence
//ff:what test: TestCollectChon
package chain

import "testing"

func TestCollectChon(t *testing.T) {
	seen := map[string]bool{"a": true}
	results := CollectChon([]string{"a", "b", "c"}, 2, "child", seen)
	if len(results) != 2 {
		t.Fatalf("len = %d, want 2", len(results))
	}
	if results[0].Name != "b" || results[0].Chon != 2 || results[0].Rel != "child" {
		t.Errorf("results[0] = %+v", results[0])
	}
}
