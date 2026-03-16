//ff:func feature=validate type=util control=sequence
//ff:what FuncDecl 내에 백틱 문자열 리터럴이 존재하는지 판별
package validate

import (
	"go/ast"
	"go/token"
	"strings"
)

// HasBacktick returns true if the func contains a backtick string literal.
func HasBacktick(fd *ast.FuncDecl) bool {
	found := false
	ast.Inspect(fd, func(n ast.Node) bool {
		lit, ok := n.(*ast.BasicLit)
		if ok && lit.Kind == token.STRING && strings.HasPrefix(lit.Value, "`") {
			found = true
			return false
		}
		return true
	})
	return found
}
