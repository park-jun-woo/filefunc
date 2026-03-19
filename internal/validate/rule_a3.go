//ff:func feature=validate type=rule control=sequence
//ff:what A3 toulmin rule — func/type 파일의 //ff:what 필수 위반 시 violation 반환
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// RuleA3 returns (true, []model.Violation) if the file violates A3 (missing what).
func RuleA3(claim any, ground any, backing any) (bool, any) {
	gf := ground.(*ValidateGround).File
	if len(gf.Funcs) == 0 && len(gf.Types) == 0 {
		return false, nil
	}
	if gf.Annotation == nil || gf.Annotation.What == "" {
		return true, []model.Violation{{
			File:    gf.Path,
			Rule:    "A3",
			Level:   "ERROR",
			Message: "file with func or type must have //ff:what annotation",
		}}
	}
	return false, nil
}
