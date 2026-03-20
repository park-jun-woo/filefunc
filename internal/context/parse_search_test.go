//ff:func feature=context type=util control=sequence
//ff:what test: TestParseSearch
package context

import "testing"

func TestParseSearch(t *testing.T) {
	got := ParseSearch("feature=validate type=rule")
	if got["feature"] != "validate" || got["type"] != "rule" {
		t.Errorf("ParseSearch = %v", got)
	}
	empty := ParseSearch("")
	if len(empty) != 0 {
		t.Errorf("ParseSearch empty = %v", empty)
	}
}
