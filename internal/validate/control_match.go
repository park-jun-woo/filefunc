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
	gf := ground.(*ValidateGround).File
	if len(gf.Funcs) == 0 || gf.Annotation == nil {
		return false, nil
	}
	if gf.Annotation.Func["control"] != b.Control {
		return false, nil
	}

	if b.MustHave != "" {
		actual := parse.DetectControl(gf.Path)
		if actual != b.Control {
			return true, []model.Violation{{
				File:    gf.Path,
				Rule:    b.Rule,
				Level:   "ERROR",
				Message: b.Message,
			}}
		}
		return false, nil
	}

	if b.MustNotHave != "" {
		return checkForbiddenControl(gf, b)
	}

	return false, nil
}
