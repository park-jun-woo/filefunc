//ff:func feature=parse type=parser
//ff:what 단일 AST 문장 노드의 nesting depth 계산
//ff:calls IfElseDepth, StmtDepth
//ff:checked llm=gpt-oss:20b hash=30f2fdaa
package parse

import "go/ast"

// NodeDepth calculates the nesting depth of a single statement node.
func NodeDepth(stmt ast.Stmt, current int) int {
	switch s := stmt.(type) {
	case *ast.IfStmt:
		return IfElseDepth(s, current)
	case *ast.ForStmt:
		return StmtDepth(s.Body.List, current+1)
	case *ast.RangeStmt:
		return StmtDepth(s.Body.List, current+1)
	case *ast.SwitchStmt:
		return StmtDepth(s.Body.List, current+1)
	case *ast.TypeSwitchStmt:
		return StmtDepth(s.Body.List, current+1)
	case *ast.SelectStmt:
		return StmtDepth(s.Body.List, current+1)
	case *ast.BlockStmt:
		return StmtDepth(s.List, current)
	case *ast.CaseClause:
		return StmtDepth(s.Body, current)
	case *ast.CommClause:
		return StmtDepth(s.Body, current)
	}
	return current
}
