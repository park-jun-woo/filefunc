//ff:func feature=parse type=parser control=sequence
//ff:what depth 1에 ForStmt/RangeStmt가 존재하는지 판별
package parse

import (
	"go/parser"
	"go/token"
)

// HasLoopAtDepth1 returns true if the first non-init func has a loop at depth 1.
func HasLoopAtDepth1(path string) bool {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return false
	}
	stmts := firstFuncBodyStmts(f)
	return containsLoop(stmts)
}
