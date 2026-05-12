//ff:func feature=validate type=util control=selection
//ff:what ExistsWhen의 존재 대상을 SourceFile에서 판별
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

func checkNeed(sf model.SourceFile, need string) bool {
	ann := sf.GetAnnotation()
	switch need {
	case "ff:func":
		return ann != nil && len(ann.Func) > 0
	case "ff:type":
		return ann != nil && len(ann.Type) > 0
	case "ff:what":
		return ann != nil && ann.What != ""
	case "control":
		if ann == nil {
			return false
		}
		c := ann.Func["control"]
		return c == "sequence" || c == "selection" || c == "iteration"
	case "dimension":
		return ann != nil && ann.Func["dimension"] != ""
	case "companion":
		return len(sf.GetFuncs()) > 0 || len(sf.GetVars()) > 0 || len(sf.GetMethods()) > 0 || len(sf.GetTypes()) > 0
	}
	return false
}
