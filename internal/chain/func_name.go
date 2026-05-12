//ff:func feature=chain type=util control=sequence
//ff:what SourceFile에서 대표 func/method 이름을 추출
package chain

import "github.com/park-jun-woo/filefunc/internal/model"

func funcName(sf model.SourceFile) string {
	if len(sf.GetFuncs()) > 0 {
		return sf.GetFuncs()[0]
	}
	if len(sf.GetMethods()) > 0 {
		return sf.GetMethods()[0]
	}
	return ""
}
