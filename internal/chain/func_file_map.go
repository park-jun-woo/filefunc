//ff:func feature=chain type=loader control=iteration dimension=1
//ff:what GoFile 목록에서 함수명/메서드명/타입명 → GoFile 매핑을 구축
package chain

import "github.com/park-jun-woo/filefunc/internal/model"

// BuildFuncFileMap creates a mapping from qualified names (pkg.Name) to their GoFile.
func BuildFuncFileMap(files []*model.GoFile) map[string]*model.GoFile {
	m := make(map[string]*model.GoFile, len(files))
	for _, gf := range files {
		for _, name := range gf.Funcs {
			m[qualifiedName(gf.Package, name)] = gf
		}
		for _, name := range gf.Methods {
			m[qualifiedName(gf.Package, name)] = gf
		}
		for _, name := range gf.Types {
			m[qualifiedName(gf.Package, name)] = gf
		}
	}
	return m
}
