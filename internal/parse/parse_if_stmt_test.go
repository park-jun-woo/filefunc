//ff:func feature=parse type=util control=sequence
//ff:what test: parseIfStmt
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
