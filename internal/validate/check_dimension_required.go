//ff:func feature=validate type=rule control=sequence
//ff:what A15: control=iteration이면 dimension= 필수
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// CheckDimensionRequired checks A15: iteration files must have dimension= annotation.
func CheckDimensionRequired(gf *model.GoFile) []model.Violation {
	if len(gf.Funcs) == 0 || gf.Annotation == nil {
		return nil
	}
	if gf.Annotation.Func["control"] != "iteration" {
		return nil
	}
	if gf.Annotation.Func["dimension"] == "" {
		return []model.Violation{{
			File:    gf.Path,
			Rule:    "A15",
			Level:   "ERROR",
			Message: "control=iteration requires dimension= annotation",
		}}
	}
	return nil
}
