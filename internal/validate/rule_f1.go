//ff:func feature=validate type=rule control=sequence
//ff:what F1 toulmin rule — 파일당 func 1개 위반 시 violation 반환
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// RuleF1 returns (true, []model.Violation) if the file violates F1 (multiple funcs).
func RuleF1(claim any, ground any) (bool, any) {
	gf := ground.(*ValidateGround).File
	if len(gf.Funcs) > 1 {
		return true, []model.Violation{{
			File:    gf.Path,
			Rule:    "F1",
			Level:   "ERROR",
			Message: "file contains multiple funcs; expected 1 file 1 func",
		}}
	}
	return false, nil
}
