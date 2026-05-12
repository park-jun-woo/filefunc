//ff:func feature=validate type=rule control=sequence
//ff:what TypeScript 전용 제어 구조 일치 검증 — TypeScriptFile 필드로 A10-A14 판정
package validate

import (
	"github.com/park-jun-woo/filefunc/internal/model"
)

// ControlMatchTypeScript returns (true, []model.Violation) if the declared control
// annotation does not match the actual TypeScript control structure at depth 1.
func ControlMatchTypeScript(claim any, ground any, backing any) (bool, any) {
	b := backing.(*ControlMatchBacking)
	sf := ground.(*ValidateGround).File
	tf, ok := sf.(*model.TypeScriptFile)
	if !ok {
		return false, nil
	}
	ann := tf.GetAnnotation()
	if ann == nil || ann.Func["control"] != b.Control {
		return false, nil
	}

	if b.MustHave != "" {
		if !hasTypeScriptControl(tf, b.MustHave) {
			return true, []model.Violation{{
				File:    tf.GetPath(),
				Rule:    b.Rule,
				Level:   "ERROR",
				Message: b.Message,
			}}
		}
		return false, nil
	}

	if b.MustNotHave != "" {
		return checkForbiddenControlTypeScript(tf, b)
	}

	return false, nil
}
