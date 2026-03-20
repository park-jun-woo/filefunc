//ff:func feature=annotate type=util control=iteration dimension=1
//ff:what test: TestReplaceAnnotationLine
package annotate

import "testing"

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
