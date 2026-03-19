//ff:func feature=validate type=rule control=sequence
//ff:what A15 toulmin rule — control=iteration이면 dimension= 필수 위반 시 violation 반환
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// RuleA15 returns (true, []model.Violation) if the file violates A15.
func RuleA15(claim any, ground any, backing any) (bool, any) {
	gf := ground.(*ValidateGround).File
	if len(gf.Funcs) == 0 || gf.Annotation == nil {
		return false, nil
	}
	if gf.Annotation.Func["control"] != "iteration" {
		return false, nil
	}
	if gf.Annotation.Func["dimension"] == "" {
		return true, []model.Violation{{
			File:    gf.Path,
			Rule:    "A15",
			Level:   "ERROR",
			Message: "control=iteration requires dimension= annotation",
		}}
	}
	return false, nil
}
