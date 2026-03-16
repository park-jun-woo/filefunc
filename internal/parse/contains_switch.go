//ff:func feature=parse type=parser control=iteration dimension=1
//ff:what statement 목록에 SwitchStmt/TypeSwitchStmt가 포함되어 있는지 판별
package parse

import "go/ast"

// containsSwitch returns true if any statement is a SwitchStmt or TypeSwitchStmt.
func containsSwitch(stmts []ast.Stmt) bool {
	for _, stmt := range stmts {
		if _, ok := stmt.(*ast.SwitchStmt); ok {
			return true
		}
		if _, ok := stmt.(*ast.TypeSwitchStmt); ok {
			return true
		}
	}
	return false
}
