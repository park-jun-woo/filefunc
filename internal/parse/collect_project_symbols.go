//ff:func feature=parse type=parser
//ff:what 프로젝트 전체 GoFile에서 func/type 이름을 수집하여 심볼 맵 반환
//ff:checked llm=gpt-oss:20b hash=201498f9
package parse

import "github.com/park-jun-woo/filefunc/internal/model"

// CollectProjectSymbols collects all func and type names from parsed GoFiles.
// Returns two maps: funcs[name]=true, types[name]=true.
func CollectProjectSymbols(files []*model.GoFile) (funcs map[string]bool, types map[string]bool) {
	funcs = make(map[string]bool)
	types = make(map[string]bool)
	for _, gf := range files {
		for _, name := range gf.Funcs {
			funcs[name] = true
		}
		for _, name := range gf.Types {
			types[name] = true
		}
	}
	return funcs, types
}
