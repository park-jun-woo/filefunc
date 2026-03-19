//ff:func feature=validate type=util control=selection
//ff:what ExistsWhen의 존재 대상을 GoFile에서 판별
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

func checkNeed(gf *model.GoFile, need string) bool {
	switch need {
	case "ff:func":
		return gf.Annotation != nil && len(gf.Annotation.Func) > 0
	case "ff:type":
		return gf.Annotation != nil && len(gf.Annotation.Type) > 0
	case "ff:what":
		return gf.Annotation != nil && gf.Annotation.What != ""
	case "control":
		if gf.Annotation == nil {
			return false
		}
		c := gf.Annotation.Func["control"]
		return c == "sequence" || c == "selection" || c == "iteration"
	case "dimension":
		return gf.Annotation != nil && gf.Annotation.Func["dimension"] != ""
	case "companion":
		return len(gf.Funcs) > 0 || len(gf.Vars) > 0 || len(gf.Methods) > 0 || len(gf.Types) > 0
	}
	return false
}
