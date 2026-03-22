//ff:func feature=validate type=util control=selection
//ff:what Q4 depth-1 제어문을 종류별로 분기하여 PURE body 검사
package validate

import (
	"go/ast"
	"go/token"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func checkQ4Stmt(fset *token.FileSet, path, funcName string, stmt ast.Stmt, violations []model.Violation) []model.Violation {
	switch s := stmt.(type) {
	case *ast.IfStmt:
		return checkQ4Block(fset, path, funcName, "if", s.Body, violations)
	case *ast.ForStmt:
		return checkQ4Block(fset, path, funcName, "for", s.Body, violations)
	case *ast.RangeStmt:
		return checkQ4Block(fset, path, funcName, "range", s.Body, violations)
	case *ast.SwitchStmt:
		return checkQ4Cases(fset, path, funcName, s.Body, violations)
	case *ast.TypeSwitchStmt:
		return checkQ4Cases(fset, path, funcName, s.Body, violations)
	}
	return violations
}
