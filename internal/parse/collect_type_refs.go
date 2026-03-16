//ff:func feature=parse type=parser control=sequence
//ff:what FuncDecl에서 프로젝트 내부 타입 참조를 수집
//ff:checked llm=gpt-oss:20b hash=28e2e40d
package parse

import "go/ast"

// CollectTypeRefs collects project-internal type references from a FuncDecl.
func CollectTypeRefs(fd *ast.FuncDecl, projImports map[string]bool, projTypes map[string]bool, seen map[string]bool) {
	ast.Inspect(fd, func(n ast.Node) bool {
		sel, ok := n.(*ast.SelectorExpr)
		if ok {
			ident, ok := sel.X.(*ast.Ident)
			if ok && projImports[ident.Name] {
				seen[sel.Sel.Name] = true
			}
			return false
		}
		ident, ok := n.(*ast.Ident)
		if ok && projTypes[ident.Name] {
			seen[ident.Name] = true
		}
		return true
	})
}
