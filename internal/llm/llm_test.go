package llm

import (
	"strings"
	"testing"
)

func TestBuildPrompt(t *testing.T) {
	p := BuildPrompt("validates file structure", "func Check() {}")
	if !strings.Contains(p, "validates file structure") {
		t.Error("missing what")
	}
	if !strings.Contains(p, "func Check() {}") {
		t.Error("missing body")
	}
	if !strings.Contains(p, "0.0") || !strings.Contains(p, "1.0") {
		t.Error("missing score range")
	}
}

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
