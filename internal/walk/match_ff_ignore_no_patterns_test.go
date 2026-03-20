//ff:func feature=parse type=util control=sequence
//ff:what test: TestMatchFFIgnore_NoPatterns
package walk

import "testing"

func TestMatchFFIgnore_NoPatterns(t *testing.T) {
	if MatchFFIgnore("src/main.go", "main.go", false, nil) {
		t.Error("expected false with no patterns")
	}
}
