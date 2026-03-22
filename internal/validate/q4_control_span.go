//ff:func feature=validate type=util control=selection
//ff:what 단일 stmt가 제어문이면 차지하는 줄수를 반환, 아니면 0
package validate

import (
	"go/ast"
	"go/token"
)

func q4ControlSpan(fset *token.FileSet, stmt ast.Stmt) int {
	switch stmt.(type) {
	case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.SwitchStmt, *ast.TypeSwitchStmt:
		start := fset.Position(stmt.Pos()).Line
		end := fset.Position(stmt.End()).Line
		return end - start + 1
	}
	return 0
}
