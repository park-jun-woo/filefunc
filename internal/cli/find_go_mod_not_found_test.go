//ff:func feature=cli type=util control=sequence
//ff:what test: TestFindGoMod_NotFound
package cli

import "testing"

func TestFindGoMod_NotFound(t *testing.T) {
	got := FindGoMod("/nonexistent/deep/path")
	if got != "go.mod" {
		t.Errorf("FindGoMod = %q, want %q", got, "go.mod")
	}
}
