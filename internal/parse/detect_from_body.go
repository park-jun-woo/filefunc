//ff:func feature=parse type=parser control=iteration dimension=1
//ff:what statement 목록에서 switch/loop 존재 여부로 제어구조 판별
package parse

import "go/ast"

func detectFromBody(stmts []ast.Stmt) string {
	for _, stmt := range stmts {
		switch stmt.(type) {
		case *ast.SwitchStmt, *ast.TypeSwitchStmt:
			return "selection"
		case *ast.ForStmt, *ast.RangeStmt:
			return "iteration"
		}
	}
	return "sequence"
}
