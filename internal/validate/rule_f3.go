//ff:func feature=validate type=rule control=sequence
//ff:what F3 toulmin rule — 파일당 method 1개 위반 시 violation 반환
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// RuleF3 returns (true, []model.Violation) if the file violates F3 (multiple methods).
func RuleF3(claim any, ground any, backing any) (bool, any) {
	gf := ground.(*ValidateGround).File
	if len(gf.Methods) > 1 {
		return true, []model.Violation{{
			File:    gf.Path,
			Rule:    "F3",
			Level:   "ERROR",
			Message: "file contains multiple methods; expected 1 file 1 method",
		}}
	}
	return false, nil
}
