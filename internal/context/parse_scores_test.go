//ff:func feature=context type=util control=iteration dimension=1
//ff:what test: TestParseScores
package context

import "testing"

func TestParseScores(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"0.85", 1},
		{"1. 0.85\n2. 0.70", 2},
		{"<think>reasoning</think>\n0.85", 1},
		{"not a score", 0},
	}
	for _, tt := range tests {
		got := ParseScores(tt.input)
		if len(got) != tt.want {
			t.Errorf("ParseScores(%q) len = %d, want %d", tt.input, len(got), tt.want)
		}
	}
}
