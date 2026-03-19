//ff:func feature=validate type=rule control=sequence
//ff:what A10 toulmin rule — control=selection인데 switch 없으면 violation 반환
package validate

import (
	"github.com/park-jun-woo/filefunc/internal/model"
	"github.com/park-jun-woo/filefunc/internal/parse"
)

// RuleA10 returns (true, []model.Violation) if the file violates A10.
func RuleA10(claim any, ground any, backing any) (bool, any) {
	gf := ground.(*ValidateGround).File
	if len(gf.Funcs) == 0 || gf.Annotation == nil {
		return false, nil
	}
	if gf.Annotation.Func["control"] != "selection" {
		return false, nil
	}
	actual := parse.DetectControl(gf.Path)
	if actual != "selection" {
		return true, []model.Violation{{
			File:    gf.Path,
			Rule:    "A10",
			Level:   "ERROR",
			Message: "control=selection but no switch found at depth 1",
		}}
	}
	return false, nil
}
