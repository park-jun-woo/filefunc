//ff:func feature=context type=util control=iteration dimension=1
//ff:what test: TestParseFeatures
package context

import "testing"

func TestParseFeatures(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{`["validate", "chain"]`, 2},
		{`some text ["validate"]`, 1},
		{`invalid`, 0},
		{`[]`, 0},
	}
	for _, tt := range tests {
		got := ParseFeatures(tt.input)
		if len(got) != tt.want {
			t.Errorf("ParseFeatures(%q) len = %d, want %d", tt.input, len(got), tt.want)
		}
	}
}
