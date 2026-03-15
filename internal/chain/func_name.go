//ff:func feature=chain type=util
//ff:what GoFile에서 대표 func/method 이름을 추출
package chain

import "github.com/park-jun-woo/filefunc/internal/model"

func funcName(gf *model.GoFile) string {
	if len(gf.Funcs) > 0 {
		return gf.Funcs[0]
	}
	if len(gf.Methods) > 0 {
		return gf.Methods[0]
	}
	return ""
}
