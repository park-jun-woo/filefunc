//ff:func feature=parse type=util control=sequence
//ff:what test: TestMatchFFIgnore_GlobPattern
package walk

import "testing"

func TestMatchFFIgnore_GlobPattern(t *testing.T) {
	patterns := []string{"*.json", "vendor/"}
	if !MatchFFIgnore("config.json", "config.json", false, patterns) {
		t.Error("expected *.json to match config.json")
	}
	if MatchFFIgnore("main.go", "main.go", false, patterns) {
		t.Error("expected *.json not to match main.go")
	}
}
