//ff:func feature=validate type=rule control=sequence
//ff:what A14 toulmin rule — control=iteration인데 switch 존재 시 violation 반환
package validate

import (
	"github.com/park-jun-woo/filefunc/internal/model"
	"github.com/park-jun-woo/filefunc/internal/parse"
)

// RuleA14 returns (true, []model.Violation) if the file violates A14.
func RuleA14(claim any, ground any) (bool, any) {
	gf := ground.(*ValidateGround).File
	if len(gf.Funcs) == 0 || gf.Annotation == nil {
		return false, nil
	}
	if gf.Annotation.Func["control"] != "iteration" {
		return false, nil
	}
	if parse.HasSwitchAtDepth1(gf.Path) {
		return true, []model.Violation{{
			File:    gf.Path,
			Rule:    "A14",
			Level:   "ERROR",
			Message: "control=iteration but switch found at depth 1; extract switch to separate func",
		}}
	}
	return false, nil
}
