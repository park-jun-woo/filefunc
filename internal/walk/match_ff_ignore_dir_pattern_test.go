//ff:func feature=parse type=util control=sequence
//ff:what test: TestMatchFFIgnore_DirPattern
package walk

import "testing"

func TestMatchFFIgnore_DirPattern(t *testing.T) {
	patterns := []string{"vendor/"}
	if !MatchFFIgnore("vendor", "vendor", true, patterns) {
		t.Error("expected vendor/ to match vendor dir")
	}
	if MatchFFIgnore("vendor", "vendor", false, patterns) {
		t.Error("expected vendor/ not to match vendor file")
	}
}
