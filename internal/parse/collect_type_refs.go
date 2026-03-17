//ff:func feature=parse type=parser control=sequence
//ff:what FuncDecl에서 프로젝트 내부 타입 참조를 수집
package parse

import "go/ast"

// CollectTypeRefs collects project-internal type references from a FuncDecl.
func CollectTypeRefs(fd *ast.FuncDecl, projImports map[string]string, projTypes map[string]bool, seen map[string]bool) {
	ast.Inspect(fd, func(n ast.Node) bool {
		sel, ok := n.(*ast.SelectorExpr)
		if ok {
			ident, ok := sel.X.(*ast.Ident)
			if ok {
				if _, exists := projImports[ident.Name]; exists {
					seen[sel.Sel.Name] = true
				}
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
