//ff:func feature=parse type=parser
//ff:what AST의 import 목록에서 프로젝트 내부 패키지 alias 맵 생성
//ff:checked llm=gpt-oss:20b hash=e377700d
package parse

import (
	"go/ast"
	"strings"
)

// BuildImportMap builds a map of package aliases that belong to the project module.
func BuildImportMap(f *ast.File, modulePath string) map[string]bool {
	m := make(map[string]bool)
	for _, imp := range f.Imports {
		path := strings.Trim(imp.Path.Value, `"`)
		if !strings.HasPrefix(path, modulePath) {
			continue
		}
		alias := ""
		if imp.Name != nil {
			alias = imp.Name.Name
		} else {
			parts := strings.Split(path, "/")
			alias = parts[len(parts)-1]
		}
		m[alias] = true
	}
	return m
}
