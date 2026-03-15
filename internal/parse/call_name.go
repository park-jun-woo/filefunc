//ff:func feature=parse type=parser
//ff:what CallExpr에서 프로젝트 내부 함수명을 추출
//ff:checked llm=gpt-oss:20b hash=5b0d8f50
package parse

import "go/ast"

// CallName extracts the function name from a CallExpr if it's a project-internal call.
// Returns empty string for external or built-in calls.
func CallName(call *ast.CallExpr, projImports map[string]bool, projFuncs map[string]bool) string {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if ok {
		ident, ok := sel.X.(*ast.Ident)
		if ok && projImports[ident.Name] {
			return sel.Sel.Name
		}
		return ""
	}
	ident, ok := call.Fun.(*ast.Ident)
	if ok && projFuncs[ident.Name] {
		return ident.Name
	}
	return ""
}
