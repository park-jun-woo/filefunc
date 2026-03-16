//ff:func feature=parse type=parser control=iteration dimension=1
//ff:what statement 목록에 ForStmt/RangeStmt가 포함되어 있는지 판별
package parse

import "go/ast"

// containsLoop returns true if any statement is a ForStmt or RangeStmt.
func containsLoop(stmts []ast.Stmt) bool {
	for _, stmt := range stmts {
		if _, ok := stmt.(*ast.ForStmt); ok {
			return true
		}
		if _, ok := stmt.(*ast.RangeStmt); ok {
			return true
		}
	}
	return false
}
