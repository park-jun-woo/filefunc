//ff:func feature=validate type=util control=iteration dimension=1
//ff:what stmt 목록에서 제어문이 차지하는 줄수 합산 (2depth 면제 계산용)
package validate

import (
	"go/ast"
	"go/token"
)

func q4InnerControlLines(fset *token.FileSet, stmts []ast.Stmt) int {
	total := 0
	for _, stmt := range stmts {
		total += q4ControlSpan(fset, stmt)
	}
	return total
}
