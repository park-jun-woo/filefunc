//ff:func feature=parse type=util control=sequence
//ff:what 단일 FuncDecl의 body를 AST 순회하여 프로젝트 내 호출을 seen에 수집
package parse

import "go/ast"

func collectCallsFromFunc(fd *ast.FuncDecl, imports map[string]string, projFuncs map[string]string, callerPkg string, seen map[string]bool) {
	ast.Inspect(fd.Body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		name := CallName(call, imports, projFuncs, callerPkg)
		if name != "" {
			seen[name] = true
		}
		return true
	})
}
