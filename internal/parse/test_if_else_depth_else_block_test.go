//ff:func feature=parse type=util control=sequence
//ff:what test: TestIfElseDepthElseBlock
package parse

import (
	"testing"
)

func TestIfElseDepthElseBlock(t *testing.T) {
	// else block should be at same depth as if body
	src := `package main
func f() {
	if true {
		x := 1
		_ = x
	} else {
		y := 2
		_ = y
	}
}`
	s := parseIfStmt(t, src)
	got := IfElseDepth(s, 0)
	// both if and else body are at depth 1
	if got != 1 {
		t.Errorf("if-else simple: got %d, want 1", got)
	}
}
