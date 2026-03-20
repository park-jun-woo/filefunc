//ff:func feature=parse type=util control=sequence
//ff:what test: TestIfElseDepthWithElse
package parse

import (
	"testing"
)

func TestIfElseDepthWithElse(t *testing.T) {
	src := `package main
func f() {
	if true {
		x := 1
		_ = x
	} else {
		if true {
			y := 2
			_ = y
		}
	}
}`
	s := parseIfStmt(t, src)
	got := IfElseDepth(s, 0)
	if got != 2 {
		t.Errorf("if-else with nested if: got %d, want 2", got)
	}
}
