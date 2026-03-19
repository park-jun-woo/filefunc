//ff:func feature=validate type=rule control=sequence
//ff:what 전제 조건 충족 시 대상이 존재하지 않으면 violation 반환
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// ExistsWhen returns (true, []model.Violation) if a precondition is met but the required target is missing.
func ExistsWhen(claim any, ground any, backing any) (bool, any) {
	b := backing.(*ExistsWhenBacking)
	gf := ground.(*ValidateGround).File

	if !checkWhen(gf, b.When) {
		return false, nil
	}
	if checkNeed(gf, b.Need) {
		return false, nil
	}
	return true, []model.Violation{{
		File:    gf.Path,
		Rule:    b.Rule,
		Level:   b.Level,
		Message: b.Message,
	}}
}
