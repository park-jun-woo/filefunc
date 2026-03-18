package parse

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func parseIfStmt(t *testing.T, src string) *ast.IfStmt {
	t.Helper()
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", src, 0)
	if err != nil {
		t.Fatal(err)
	}
	fn := f.Decls[0].(*ast.FuncDecl)
	return fn.Body.List[0].(*ast.IfStmt)
}

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
