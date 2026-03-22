//ff:func feature=validate type=util control=sequence
//ff:what CaseClause의 body 줄수를 계산 (Colon~마지막 stmt 사이)
package validate

import (
	"go/ast"
	"go/token"
)

func q4CaseBodyLines(fset *token.FileSet, cc *ast.CaseClause) int {
	if len(cc.Body) == 0 {
		return 0
	}
	start := fset.Position(cc.Colon).Line
	last := cc.Body[len(cc.Body)-1]
	end := fset.Position(last.End()).Line
	lines := end - start
	if lines < 0 {
		return 0
	}
	return lines
}
