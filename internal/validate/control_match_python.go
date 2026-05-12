//ff:func feature=validate type=rule control=sequence
//ff:what Python 전용 제어 구조 일치 검증 — PythonFile 필드로 A10-A14 판정
package validate

import (
	"github.com/park-jun-woo/filefunc/internal/model"
)

// ControlMatchPython returns (true, []model.Violation) if the declared control
// annotation does not match the actual Python control structure at depth 1.
func ControlMatchPython(claim any, ground any, backing any) (bool, any) {
	b := backing.(*ControlMatchBacking)
	sf := ground.(*ValidateGround).File
	pf, ok := sf.(*model.PythonFile)
	if !ok {
		return false, nil
	}
	ann := pf.GetAnnotation()
	if ann == nil || ann.Func["control"] != b.Control {
		return false, nil
	}

	if b.MustHave != "" {
		if !hasPythonControl(pf, b.MustHave) {
			return true, []model.Violation{{
				File:    pf.GetPath(),
				Rule:    b.Rule,
				Level:   "ERROR",
				Message: b.Message,
			}}
		}
		return false, nil
	}

	if b.MustNotHave != "" {
		return checkForbiddenControlPython(pf, b)
	}

	return false, nil
}
