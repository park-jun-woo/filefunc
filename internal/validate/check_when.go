//ff:func feature=validate type=util control=selection
//ff:what ExistsWhen의 전제 조건을 GoFile에서 판별
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

func checkWhen(gf *model.GoFile, when string) bool {
	switch when {
	case "HasFuncs":
		return len(gf.Funcs) > 0
	case "HasTypes":
		return len(gf.Types) > 0
	case "HasFuncOrType":
		return len(gf.Funcs) > 0 || len(gf.Types) > 0
	case "HasInit":
		return gf.HasInit
	case "ControlIteration":
		return gf.Annotation != nil && gf.Annotation.Func["control"] == "iteration"
	}
	return false
}
