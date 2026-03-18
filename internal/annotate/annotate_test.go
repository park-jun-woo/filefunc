package annotate

import (
	"testing"
)

func TestReplaceAnnotationLine(t *testing.T) {
	tests := []struct {
		name    string
		trimmed string
		prefix  string
		key     string
		newLine string
		value   string
		wantRep string
		wantOK  bool
	}{
		{"match replace", "//ff:what old text", "//ff:what ", "what", "//ff:what new text", "new text", "//ff:what new text", true},
		{"match remove", "//ff:what old text", "//ff:what ", "what", "", "", "", true},
		{"exact key only", "//ff:what", "//ff:what ", "what", "//ff:what new", "new", "//ff:what new", true},
		{"no match", "//ff:func feature=x", "//ff:what ", "what", "//ff:what x", "x", "//ff:func feature=x", false},
		{"non-annotation", "package main", "//ff:what ", "what", "//ff:what x", "x", "package main", false},
	}
	for _, tt := range tests {
		rep, ok := ReplaceAnnotationLine(tt.trimmed, tt.prefix, tt.key, tt.newLine, tt.value)
		if rep != tt.wantRep || ok != tt.wantOK {
			t.Errorf("%s: got (%q, %v), want (%q, %v)", tt.name, rep, ok, tt.wantRep, tt.wantOK)
		}
	}
}

func TestInsertAfterAnnotations(t *testing.T) {
	lines := []string{"//ff:func feature=x", "//ff:what hello", "package main", ""}
	result := InsertAfterAnnotations(lines, "//ff:why reason")
	if len(result) != 5 {
		t.Fatalf("len = %d, want 5", len(result))
	}
	if result[2] != "//ff:why reason" {
		t.Errorf("result[2] = %q, want %q", result[2], "//ff:why reason")
	}
	if result[3] != "package main" {
		t.Errorf("result[3] = %q, want %q", result[3], "package main")
	}
}

func TestInsertAfterAnnotationsNoAnnotations(t *testing.T) {
	lines := []string{"package main", "func foo() {}"}
	result := InsertAfterAnnotations(lines, "//ff:what new")
	if len(result) != 3 {
		t.Fatalf("len = %d, want 3", len(result))
	}
	if result[0] != "//ff:what new" {
		t.Errorf("result[0] = %q, want %q", result[0], "//ff:what new")
	}
}
