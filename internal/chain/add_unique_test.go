//ff:func feature=chain type=util control=sequence
//ff:what test: TestAddUnique
package chain

import "testing"

func TestAddUnique(t *testing.T) {
	seen := map[string]bool{"a": true}
	var result []string
	AddUnique([]string{"a", "b", "c", "b"}, seen, &result)
	if len(result) != 2 {
		t.Fatalf("len = %d, want 2", len(result))
	}
	if result[0] != "b" || result[1] != "c" {
		t.Errorf("result = %v, want [b c]", result)
	}
}
