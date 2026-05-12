//ff:func feature=validate type=rule control=sequence
//ff:what 선언된 제어 구조와 AST 실체가 불일치하면 violation 반환
package validate

import (
	"github.com/park-jun-woo/filefunc/internal/model"
	"github.com/park-jun-woo/filefunc/internal/parse"
)

// ControlMatch returns (true, []model.Violation) if the declared control annotation
// does not match the actual AST control structure at depth 1.
func ControlMatch(claim any, ground any, backing any) (bool, any) {
	b := backing.(*ControlMatchBacking)
	sf := ground.(*ValidateGround).File
	ann := sf.GetAnnotation()
	if len(sf.GetFuncs()) == 0 || ann == nil {
		return false, nil
	}
	if ann.Func["control"] != b.Control {
		return false, nil
	}

	if b.MustHave != "" {
		actual := parse.DetectControl(sf.GetPath())
		if actual != b.Control {
			return true, []model.Violation{{
				File:    sf.GetPath(),
				Rule:    b.Rule,
				Level:   "ERROR",
				Message: b.Message,
			}}
		}
		return false, nil
	}

	if b.MustNotHave != "" {
		return checkForbiddenControl(sf, b)
	}

	return false, nil
}
