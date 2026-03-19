//ff:func feature=validate type=rule control=sequence
//ff:what F4 toulmin rule — init()만 단독 존재하는 파일 위반 시 violation 반환
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// RuleF4 returns (true, []model.Violation) if the file violates F4 (init-only file).
func RuleF4(claim any, ground any, backing any) (bool, any) {
	gf := ground.(*ValidateGround).File
	if !gf.HasInit {
		return false, nil
	}
	if len(gf.Funcs) == 0 && len(gf.Vars) == 0 && len(gf.Methods) == 0 && len(gf.Types) == 0 {
		return true, []model.Violation{{
			File:    gf.Path,
			Rule:    "F4",
			Level:   "ERROR",
			Message: "init() must not exist alone; requires accompanying var or func",
		}}
	}
	return false, nil
}
