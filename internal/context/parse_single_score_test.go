//ff:func feature=context type=util control=iteration dimension=1
//ff:what test: TestParseSingleScore
package context

import "testing"

func TestParseSingleScore(t *testing.T) {
	tests := []struct {
		input string
		want  float64
	}{
		{"0.85", 0.85},
		{"1. 0.70", 0.70},
		{"abc", -1},
		{"", -1},
	}
	for _, tt := range tests {
		got := parseSingleScore(tt.input)
		if got != tt.want {
			t.Errorf("parseSingleScore(%q) = %f, want %f", tt.input, got, tt.want)
		}
	}
}
