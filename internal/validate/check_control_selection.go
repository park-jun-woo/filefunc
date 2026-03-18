//ff:func feature=validate type=rule control=sequence
//ff:what A10: control=selection인데 depth 1에 switch 없으면 ERROR
package validate

import (
	"github.com/park-jun-woo/filefunc/internal/model"
	"github.com/park-jun-woo/filefunc/internal/parse"
)

// CheckControlSelection checks A10: control=selection requires switch at depth 1.
func CheckControlSelection(gf *model.GoFile) []model.Violation {
	if len(gf.Funcs) == 0 || gf.Annotation == nil {
		return nil
	}
	if gf.Annotation.Func["control"] != "selection" {
		return nil
	}
	actual := parse.DetectControl(gf.Path)
	if actual != "selection" {
		return []model.Violation{{
			File:    gf.Path,
			Rule:    "A10",
			Level:   "ERROR",
			Message: "control=selection but no switch found at depth 1",
		}}
	}
	return nil
}
