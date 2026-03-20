//ff:func feature=chain type=util control=iteration dimension=1
//ff:what test: TestPkgFromQualified
package chain

import "testing"

func TestPkgFromQualified(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"pkg.FuncName", "pkg"},
		{"a.b.c.Func", "a.b.c"},
		{"NoPackage", ""},
		{"", ""},
	}
	for _, tt := range tests {
		got := PkgFromQualified(tt.input)
		if got != tt.want {
			t.Errorf("PkgFromQualified(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
