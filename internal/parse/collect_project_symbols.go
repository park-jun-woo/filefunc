//ff:func feature=parse type=parser control=iteration dimension=1
//ff:what 프로젝트 전체 SourceFile에서 func/type 이름을 수집하여 심볼 맵 반환
package parse

import "github.com/park-jun-woo/filefunc/internal/model"

// CollectProjectSymbols collects all func and type names from parsed SourceFiles.
// Returns two maps: funcs[name]=pkg, types[name]=true.
func CollectProjectSymbols(files []model.SourceFile) (funcs map[string]string, types map[string]bool) {
	funcs = make(map[string]string)
	types = make(map[string]bool)
	for _, sf := range files {
		for _, name := range sf.GetFuncs() {
			funcs[name] = sf.GetPackage()
		}
		for _, name := range sf.GetMethods() {
			funcs[name] = sf.GetPackage()
		}
		for _, name := range sf.GetTypes() {
			types[name] = true
		}
	}
	return funcs, types
}
