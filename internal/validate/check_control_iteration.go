//ff:func feature=validate type=rule control=sequence
//ff:what A10: control=iteration인데 depth 1에 loop 없으면 ERROR
package validate

import (
	"github.com/park-jun-woo/filefunc/internal/model"
	"github.com/park-jun-woo/filefunc/internal/parse"
)

// CheckControlIteration checks A10: control=iteration requires loop at depth 1.
func CheckControlIteration(gf *model.GoFile) []model.Violation {
	if gf.IsTest || len(gf.Funcs) == 0 || gf.Annotation == nil {
		return nil
	}
	if gf.Annotation.Func["control"] != "iteration" {
		return nil
	}
	actual := parse.DetectControl(gf.Path)
	if actual != "iteration" {
		return []model.Violation{{
			File:    gf.Path,
			Rule:    "A10",
			Level:   "ERROR",
			Message: "control=iteration but no loop found at depth 1",
		}}
	}
	return nil
}
