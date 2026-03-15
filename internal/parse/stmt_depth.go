//ff:func feature=parse type=parser
//ff:what 문장 목록의 최대 nesting depth 계산
package parse

import "go/ast"

// StmtDepth calculates the maximum nesting depth in a list of statements.
func StmtDepth(stmts []ast.Stmt, current int) int {
	maxDepth := current
	for _, stmt := range stmts {
		d := NodeDepth(stmt, current)
		if d > maxDepth {
			maxDepth = d
		}
	}
	return maxDepth
}
