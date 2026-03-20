//ff:func feature=parse type=util control=sequence
//ff:what test: TestIfElseDepthSimple
package parse

import (
	"testing"
)

func TestIfElseDepthSimple(t *testing.T) {
	src := `package main
func f() {
	if true {
		x := 1
		_ = x
	}
}`
	s := parseIfStmt(t, src)
	got := IfElseDepth(s, 0)
	if got != 1 {
		t.Errorf("simple if: got %d, want 1", got)
	}
}
