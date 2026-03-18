//ff:func feature=validate type=rule control=sequence
//ff:what A9: func 파일은 control= 어노테이션 필수 (sequence/selection/iteration)
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// CheckControlRequired checks A9: func files must have control= annotation
// with a valid value (sequence, selection, iteration).
func CheckControlRequired(gf *model.GoFile) []model.Violation {
	if len(gf.Funcs) == 0 || gf.Annotation == nil {
		return nil
	}
	control := gf.Annotation.Func["control"]
	if control == "sequence" || control == "selection" || control == "iteration" {
		return nil
	}
	if control == "" {
		return []model.Violation{{
			File:    gf.Path,
			Rule:    "A9",
			Level:   "ERROR",
			Message: "func file must have control= annotation (sequence, selection, or iteration)",
		}}
	}
	return []model.Violation{{
		File:    gf.Path,
		Rule:    "A9",
		Level:   "ERROR",
		Message: "invalid control value: must be sequence, selection, or iteration",
	}}
}
