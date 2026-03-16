//ff:func feature=parse type=parser control=iteration dimension=1
//ff:what AST에서 첫 번째 non-init FuncDecl의 body statement 목록을 반환
package parse

import "go/ast"

// firstFuncBodyStmts returns the body statements of the first non-init FuncDecl.
func firstFuncBodyStmts(f *ast.File) []ast.Stmt {
	for _, decl := range f.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok || fd.Body == nil || fd.Name.Name == "init" {
			continue
		}
		return fd.Body.List
	}
	return nil
}
