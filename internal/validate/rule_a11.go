//ff:func feature=validate type=rule control=sequence
//ff:what A11 toulmin rule — control=iteration인데 loop 없으면 violation 반환
package validate

import (
	"github.com/park-jun-woo/filefunc/internal/model"
	"github.com/park-jun-woo/filefunc/internal/parse"
)

// RuleA11 returns (true, []model.Violation) if the file violates A11.
func RuleA11(claim any, ground any, backing any) (bool, any) {
	gf := ground.(*ValidateGround).File
	if len(gf.Funcs) == 0 || gf.Annotation == nil {
		return false, nil
	}
	if gf.Annotation.Func["control"] != "iteration" {
		return false, nil
	}
	actual := parse.DetectControl(gf.Path)
	if actual != "iteration" {
		return true, []model.Violation{{
			File:    gf.Path,
			Rule:    "A11",
			Level:   "ERROR",
			Message: "control=iteration but no loop found at depth 1",
		}}
	}
	return false, nil
}
