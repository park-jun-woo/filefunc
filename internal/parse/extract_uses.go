//ff:func feature=parse type=parser control=iteration dimension=1
//ff:what func에서 프로젝트 내 사용 타입명을 AST로 추출
//ff:checked llm=gpt-oss:20b hash=085874a2
package parse

import (
	"go/ast"
	"go/parser"
	"go/token"
)

// ExtractUses extracts project-internal type usages from a Go source file.
func ExtractUses(path string, modulePath string, projTypes map[string]bool) ([]string, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil, err
	}

	imports := BuildImportMap(f, modulePath)
	seen := make(map[string]bool)

	for _, decl := range f.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok || fd.Name.Name == "init" {
			continue
		}
		CollectTypeRefs(fd, imports, projTypes, seen)
	}

	return SortedKeys(seen), nil
}
