//ff:func feature=validate type=util control=selection
//ff:what ExistsWhen의 전제 조건을 SourceFile에서 판별
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

func checkWhen(sf model.SourceFile, when string) bool {
	switch when {
	case "HasFuncs":
		return len(sf.GetFuncs()) > 0
	case "HasTypes":
		return len(sf.GetTypes()) > 0
	case "HasFuncOrType":
		return len(sf.GetFuncs()) > 0 || len(sf.GetTypes()) > 0
	case "HasInit":
		return sf.GetHasInit()
	case "ControlIteration":
		ann := sf.GetAnnotation()
		return ann != nil && ann.Func["control"] == "iteration"
	}
	return false
}
