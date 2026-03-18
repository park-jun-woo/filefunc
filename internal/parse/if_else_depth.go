//ff:func feature=parse type=parser control=selection
//ff:what if문의 else 분기를 포함한 최대 nesting depth 계산
package parse

import "go/ast"

// IfElseDepth calculates the maximum nesting depth for an if statement,
// including its else branch.
func IfElseDepth(s *ast.IfStmt, current int) int {
	d := StmtDepth(s.Body.List, current+1)
	if s.Else == nil {
		return d
	}
	var ed int
	switch e := s.Else.(type) {
	case *ast.IfStmt:
		ed = IfElseDepth(e, current)
	case *ast.BlockStmt:
		ed = StmtDepth(e.List, current+1)
	default:
		ed = NodeDepth(s.Else, current)
	}
	if ed > d {
		return ed
	}
	return d
}
