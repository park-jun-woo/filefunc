//ff:func feature=validate type=rule control=sequence
//ff:what F2 toulmin rule — 파일당 type 1개 위반 시 violation 반환
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// RuleF2 returns (true, []model.Violation) if the file violates F2 (multiple types).
func RuleF2(claim any, ground any, backing any) (bool, any) {
	gf := ground.(*ValidateGround).File
	if len(gf.Types) > 1 {
		return true, []model.Violation{{
			File:    gf.Path,
			Rule:    "F2",
			Level:   "ERROR",
			Message: "file contains multiple types; expected 1 file 1 type",
		}}
	}
	return false, nil
}
