//ff:func feature=parse type=parser
//ff:what if문의 else 분기를 포함한 최대 nesting depth 계산
//ff:checked llm=gpt-oss:20b hash=02db6e62
package parse

import "go/ast"

// IfElseDepth calculates the maximum nesting depth for an if statement,
// including its else branch.
func IfElseDepth(s *ast.IfStmt, current int) int {
	d := StmtDepth(s.Body.List, current+1)
	if s.Else == nil {
		return d
	}
	ed := NodeDepth(s.Else, current)
	if ed > d {
		return ed
	}
	return d
}
