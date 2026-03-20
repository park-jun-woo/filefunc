//ff:func feature=chain type=util control=sequence
//ff:what test: TestCountChon2Plus
package chain

import "testing"

func TestCountChon2Plus(t *testing.T) {
	results := []ChonResult{
		{"a", 1, "child"},
		{"b", 2, "grandchild"},
		{"c", 3, "grandchild"},
	}
	got := countChon2Plus(results)
	if got != 2 {
		t.Errorf("countChon2Plus = %d, want 2", got)
	}
}
