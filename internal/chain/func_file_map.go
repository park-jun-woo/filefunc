//ff:func feature=chain type=loader control=iteration dimension=1
//ff:what SourceFile 목록에서 함수명/메서드명/타입명 → SourceFile 매핑을 구축
package chain

import "github.com/park-jun-woo/filefunc/internal/model"

// BuildFuncFileMap creates a mapping from qualified names (pkg.Name) to their SourceFile.
func BuildFuncFileMap(files []model.SourceFile) map[string]model.SourceFile {
	m := make(map[string]model.SourceFile, len(files))
	for _, sf := range files {
		for _, name := range sf.GetFuncs() {
			m[qualifiedName(sf.GetPackage(), name)] = sf
		}
		for _, name := range sf.GetMethods() {
			m[qualifiedName(sf.GetPackage(), name)] = sf
		}
		for _, name := range sf.GetTypes() {
			m[qualifiedName(sf.GetPackage(), name)] = sf
		}
	}
	return m
}
