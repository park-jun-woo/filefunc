//ff:func feature=chain type=util control=iteration dimension=1
//ff:what test: TestNameFromQualified
package chain

import "testing"

func TestNameFromQualified(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"pkg.FuncName", "FuncName"},
		{"a.b.c.Func", "Func"},
		{"NoPackage", "NoPackage"},
		{"", ""},
	}
	for _, tt := range tests {
		got := NameFromQualified(tt.input)
		if got != tt.want {
			t.Errorf("NameFromQualified(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
