//ff:func feature=parse type=parser control=iteration dimension=1
//ff:what func body의 직계 자식(depth 1)을 순회하여 제어구조를 판별
package parse

import (
	"go/ast"
	"go/parser"
	"go/token"
)

// DetectControl detects the control structure of a func.
// Returns "selection", "iteration", or "sequence".
func DetectControl(path string) string {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return "sequence"
	}

	for _, decl := range f.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok || fd.Body == nil || fd.Name.Name == "init" {
			continue
		}
		return detectFromBody(fd.Body.List)
	}
	return "sequence"
}
