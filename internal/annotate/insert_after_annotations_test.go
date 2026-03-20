//ff:func feature=annotate type=util control=sequence
//ff:what test: TestInsertAfterAnnotations
package annotate

import "testing"

func TestInsertAfterAnnotations(t *testing.T) {
	lines := []string{"//ff:func feature=x", "//ff:what hello", "package main", ""}
	result := InsertAfterAnnotations(lines, "//ff:why reason")
	if len(result) != 5 {
		t.Fatalf("len = %d, want 5", len(result))
	}
	if result[2] != "//ff:why reason" {
		t.Errorf("result[2] = %q, want %q", result[2], "//ff:why reason")
	}
	if result[3] != "package main" {
		t.Errorf("result[3] = %q, want %q", result[3], "package main")
	}
}
