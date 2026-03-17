//ff:func feature=parse type=parser control=sequence
//ff:what CallExpr에서 프로젝트 내부 qualified 함수명(pkg.FuncName)을 추출
package parse

import "go/ast"

// CallName extracts the qualified function name (pkg.FuncName) from a CallExpr.
// projImports: alias → package name. projFuncs: funcName → package name.
// callerPkg is used for same-package calls. Returns empty string for external or built-in calls.
func CallName(call *ast.CallExpr, projImports map[string]string, projFuncs map[string]string, callerPkg string) string {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if ok {
		ident, ok := sel.X.(*ast.Ident)
		if !ok {
			return ""
		}
		pkg, exists := projImports[ident.Name]
		if exists {
			return pkg + "." + sel.Sel.Name
		}
		return ""
	}
	ident, ok := call.Fun.(*ast.Ident)
	if !ok {
		return ""
	}
	if _, exists := projFuncs[ident.Name]; exists {
		return callerPkg + "." + ident.Name
	}
	return ""
}
