//ff:func feature=cli type=util control=iteration dimension=1
//ff:what test: TestParseScore
package llm

import "testing"

func TestParseScore(t *testing.T) {
	tests := []struct {
		input   string
		want    float64
		wantErr bool
	}{
		{"0.85", 0.85, false},
		{"  1.0  ", 1.0, false},
		{"0.0", 0.0, false},
		{"1.5", 0, true},
		{"-0.1", 0, true},
		{"abc", 0, true},
	}
	for _, tt := range tests {
		got, err := ParseScore(tt.input)
		if tt.wantErr && err == nil {
			t.Errorf("ParseScore(%q) expected error", tt.input)
		}
		if !tt.wantErr && err != nil {
			t.Errorf("ParseScore(%q) unexpected error: %v", tt.input, err)
		}
		if !tt.wantErr && got != tt.want {
			t.Errorf("ParseScore(%q) = %f, want %f", tt.input, got, tt.want)
		}
	}
}
