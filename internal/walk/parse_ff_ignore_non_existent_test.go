//ff:func feature=parse type=util control=sequence
//ff:what test: TestParseFFIgnore_NonExistent
package walk

import "testing"

func TestParseFFIgnore_NonExistent(t *testing.T) {
	patterns := ParseFFIgnore("/nonexistent/.ffignore")
	if patterns != nil {
		t.Errorf("expected nil, got %v", patterns)
	}
}
