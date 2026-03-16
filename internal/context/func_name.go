//ff:func feature=context type=util control=sequence
//ff:what GoFile에서 대표 함수명/타입명을 추출
package context

import "github.com/park-jun-woo/filefunc/internal/model"

// funcName returns the primary func, method, or type name from a GoFile.
func funcName(gf *model.GoFile) string {
	if len(gf.Funcs) > 0 {
		return gf.Funcs[0]
	}
	if len(gf.Methods) > 0 {
		return gf.Methods[0]
	}
	if len(gf.Types) > 0 {
		return gf.Types[0]
	}
	return gf.Path
}
