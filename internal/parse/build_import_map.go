//ff:func feature=parse type=parser control=iteration dimension=1
//ff:what AST의 import 목록에서 프로젝트 내부 패키지 alias → 패키지명 맵 생성
package parse

import (
	"go/ast"
	"strings"
)

// BuildImportMap builds a map of alias → package name for project-internal imports.
func BuildImportMap(f *ast.File, modulePath string) map[string]string {
	m := make(map[string]string)
	for _, imp := range f.Imports {
		path := strings.Trim(imp.Path.Value, `"`)
		if !strings.HasPrefix(path, modulePath) {
			continue
		}
		parts := strings.Split(path, "/")
		pkgName := parts[len(parts)-1]
		alias := pkgName
		if imp.Name != nil {
			alias = imp.Name.Name
		}
		m[alias] = pkgName
	}
	return m
}
