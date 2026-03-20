//ff:func feature=annotate type=util control=sequence
//ff:what test: TestInsertAfterAnnotationsNoAnnotations
package annotate

import "testing"

func TestInsertAfterAnnotationsNoAnnotations(t *testing.T) {
	lines := []string{"package main", "func foo() {}"}
	result := InsertAfterAnnotations(lines, "//ff:what new")
	if len(result) != 3 {
		t.Fatalf("len = %d, want 3", len(result))
	}
	if result[0] != "//ff:what new" {
		t.Errorf("result[0] = %q, want %q", result[0], "//ff:what new")
	}
}
