//ff:func feature=parse type=parser control=iteration
//ff:what func body에서 프로젝트 내 호출 함수명을 AST로 추출
//ff:checked llm=gpt-oss:20b hash=b668ae9c
package parse

import (
	"go/ast"
	"go/parser"
	"go/token"
)

// ExtractCalls extracts project-internal function calls from a Go source file.
func ExtractCalls(path string, modulePath string, projFuncs map[string]bool) ([]string, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil, err
	}

	imports := BuildImportMap(f, modulePath)
	seen := make(map[string]bool)

	for _, decl := range f.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok || fd.Body == nil || fd.Name.Name == "init" {
			continue
		}
		ast.Inspect(fd.Body, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			name := CallName(call, imports, projFuncs)
			if name != "" {
				seen[name] = true
			}
			return true
		})
	}

	return SortedKeys(seen), nil
}
