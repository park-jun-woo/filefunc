//ff:func feature=parse type=parser control=sequence
//ff:what depth 1에 SwitchStmt/TypeSwitchStmt가 존재하는지 판별
package parse

import (
	"go/parser"
	"go/token"
)

// HasSwitchAtDepth1 returns true if the first non-init func has a switch at depth 1.
func HasSwitchAtDepth1(path string) bool {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return false
	}
	stmts := firstFuncBodyStmts(f)
	return containsSwitch(stmts)
}
