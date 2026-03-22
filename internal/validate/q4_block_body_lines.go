//ff:func feature=validate type=util control=sequence
//ff:what BlockStmt의 body 줄수를 계산 (Lbrace~Rbrace 사이)
package validate

import (
	"go/ast"
	"go/token"
)

func q4BlockBodyLines(fset *token.FileSet, block *ast.BlockStmt) int {
	start := fset.Position(block.Lbrace).Line
	end := fset.Position(block.Rbrace).Line
	lines := end - start - 1
	if lines < 0 {
		return 0
	}
	return lines
}
