//ff:func feature=validate type=rule
//ff:what C2: codebook 내 동일 키에서 중복 값을 검출
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// CheckCodebookDuplicates checks C2: no duplicate values within the same key.
func CheckCodebookDuplicates(cb *model.Codebook) []model.Violation {
	var violations []model.Violation
	violations = append(violations, FindDuplicates("feature", cb.Feature)...)
	violations = append(violations, FindDuplicates("type", cb.Type)...)
	violations = append(violations, FindDuplicates("pattern", cb.Pattern)...)
	violations = append(violations, FindDuplicates("level", cb.Level)...)
	return violations
}
