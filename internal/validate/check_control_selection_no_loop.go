//ff:func feature=validate type=rule control=sequence
//ff:what A13: control=selection인데 depth 1에 loop 존재하면 ERROR
package validate

import (
	"github.com/park-jun-woo/filefunc/internal/model"
	"github.com/park-jun-woo/filefunc/internal/parse"
)

// CheckControlSelectionNoLoop checks A13: control=selection must not have loop at depth 1.
func CheckControlSelectionNoLoop(gf *model.GoFile) []model.Violation {
	if gf.IsTest || len(gf.Funcs) == 0 || gf.Annotation == nil {
		return nil
	}
	if gf.Annotation.Func["control"] != "selection" {
		return nil
	}
	if parse.HasLoopAtDepth1(gf.Path) {
		return []model.Violation{{
			File:    gf.Path,
			Rule:    "A13",
			Level:   "ERROR",
			Message: "control=selection but loop found at depth 1; extract loop to separate func",
		}}
	}
	return nil
}
