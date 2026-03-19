//ff:func feature=validate type=rule control=sequence
//ff:what 파일 내 항목 수가 상한을 초과하면 violation 반환
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// CountMax returns (true, []model.Violation) if the count of the specified field exceeds Max.
func CountMax(claim any, ground any, backing any) (bool, any) {
	b := backing.(*CountMaxBacking)
	gf := ground.(*ValidateGround).File
	if countField(gf, b.Field) > b.Max {
		return true, []model.Violation{{
			File:    gf.Path,
			Rule:    b.Rule,
			Level:   "ERROR",
			Message: b.Message,
		}}
	}
	return false, nil
}
