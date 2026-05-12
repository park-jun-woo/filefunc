//ff:func feature=parse type=util control=iteration dimension=1
//ff:what test: TestIsSkippableLineTypeScript
package parse

import "testing"

func TestIsSkippableLineTypeScript(t *testing.T) {
	cases := []struct {
		line string
		want bool
	}{
		{"", true},
		{"// regular comment", true},
		{"// another comment", true},
		{"//ff:func feature=validate", false},
		{"import { foo } from './bar'", false},
		{"export function foo() {}", false},
		{"const x = 1", false},
	}
	for _, c := range cases {
		got := IsSkippableLineTypeScript(c.line)
		if got != c.want {
			t.Errorf("IsSkippableLineTypeScript(%q) = %v, want %v", c.line, got, c.want)
		}
	}
}
