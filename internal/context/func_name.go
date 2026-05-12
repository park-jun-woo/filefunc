//ff:func feature=context type=util control=sequence
//ff:what SourceFile에서 대표 함수명/타입명을 추출
package context

import "github.com/park-jun-woo/filefunc/internal/model"

// funcName returns the primary func, method, or type name from a SourceFile.
func funcName(sf model.SourceFile) string {
	if len(sf.GetFuncs()) > 0 {
		return sf.GetFuncs()[0]
	}
	if len(sf.GetMethods()) > 0 {
		return sf.GetMethods()[0]
	}
	if len(sf.GetTypes()) > 0 {
		return sf.GetTypes()[0]
	}
	return sf.GetPath()
}
