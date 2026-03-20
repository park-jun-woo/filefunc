//ff:func feature=parse type=util control=sequence
//ff:what test: TestIfElseIfDepth
package parse

import (
	"testing"
)

func TestIfElseIfDepth(t *testing.T) {
	src := `package main
func f() {
	if true {
		x := 1
		_ = x
	} else if false {
		y := 2
		_ = y
	}
}`
	s := parseIfStmt(t, src)
	got := IfElseDepth(s, 0)
	// else-if is at same level as if, both bodies at depth 1
	if got != 1 {
		t.Errorf("if-else-if: got %d, want 1", got)
	}
}
