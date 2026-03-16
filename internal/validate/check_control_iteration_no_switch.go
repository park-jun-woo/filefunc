//ff:func feature=validate type=rule control=sequence
//ff:what A14: control=iteration인데 depth 1에 switch 존재하면 ERROR
package validate

import (
	"github.com/park-jun-woo/filefunc/internal/model"
	"github.com/park-jun-woo/filefunc/internal/parse"
)

// CheckControlIterationNoSwitch checks A14: control=iteration must not have switch at depth 1.
func CheckControlIterationNoSwitch(gf *model.GoFile) []model.Violation {
	if gf.IsTest || len(gf.Funcs) == 0 || gf.Annotation == nil {
		return nil
	}
	if gf.Annotation.Func["control"] != "iteration" {
		return nil
	}
	if parse.HasSwitchAtDepth1(gf.Path) {
		return []model.Violation{{
			File:    gf.Path,
			Rule:    "A14",
			Level:   "ERROR",
			Message: "control=iteration but switch found at depth 1; extract switch to separate func",
		}}
	}
	return nil
}
