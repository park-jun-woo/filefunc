//ff:func feature=validate type=rule control=iteration dimension=1
//ff:what Q4 toulmin rule — 제어문 body PURE 줄수 10줄 초과 시 violation 반환
package validate

import (
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// CheckControlBody returns (true, []model.Violation) if a depth-1 control statement's
// PURE body exceeds 10 lines. PURE = total body lines minus inner control statement lines.
// For switch/type-switch, each case clause is checked individually.
func CheckControlBody(claim any, ground any, backing any) (bool, any) {
	gf := ground.(*ValidateGround).File

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, gf.Path, nil, 0)
	if err != nil {
		return false, nil
	}

	var violations []model.Violation

	for _, decl := range f.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok || fd.Body == nil {
			continue
		}
		for _, stmt := range fd.Body.List {
			violations = checkQ4Stmt(fset, gf.Path, fd.Name.Name, stmt, violations)
		}
	}

	if len(violations) > 0 {
		return true, violations
	}
	return false, nil
}
