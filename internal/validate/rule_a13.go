//ff:func feature=validate type=rule control=sequence
//ff:what A13 toulmin rule — control=selection인데 loop 존재 시 violation 반환
package validate

import (
	"github.com/park-jun-woo/filefunc/internal/model"
	"github.com/park-jun-woo/filefunc/internal/parse"
)

// RuleA13 returns (true, []model.Violation) if the file violates A13.
func RuleA13(claim any, ground any, backing any) (bool, any) {
	gf := ground.(*ValidateGround).File
	if len(gf.Funcs) == 0 || gf.Annotation == nil {
		return false, nil
	}
	if gf.Annotation.Func["control"] != "selection" {
		return false, nil
	}
	if parse.HasLoopAtDepth1(gf.Path) {
		return true, []model.Violation{{
			File:    gf.Path,
			Rule:    "A13",
			Level:   "ERROR",
			Message: "control=selection but loop found at depth 1; extract loop to separate func",
		}}
	}
	return false, nil
}
