//ff:func feature=validate type=rule control=sequence
//ff:what A9 toulmin rule — func 파일의 control= 필수 위반 시 violation 반환
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// RuleA9 returns (true, []model.Violation) if the file violates A9 (missing control).
func RuleA9(claim any, ground any, backing any) (bool, any) {
	gf := ground.(*ValidateGround).File
	if len(gf.Funcs) == 0 || gf.Annotation == nil {
		return false, nil
	}
	control := gf.Annotation.Func["control"]
	if control == "sequence" || control == "selection" || control == "iteration" {
		return false, nil
	}
	if control == "" {
		return true, []model.Violation{{
			File:    gf.Path,
			Rule:    "A9",
			Level:   "ERROR",
			Message: "func file must have control= annotation (sequence, selection, or iteration)",
		}}
	}
	return true, []model.Violation{{
		File:    gf.Path,
		Rule:    "A9",
		Level:   "ERROR",
		Message: "invalid control value: must be sequence, selection, or iteration",
	}}
}
