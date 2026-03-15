//ff:func feature=validate type=rule
//ff:what C2: codebook required+optional 전체에서 동일 키 내 중복 값을 검출
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// CheckCodebookDuplicates checks C2: no duplicate values within the same key.
func CheckCodebookDuplicates(cb *model.Codebook) []model.Violation {
	var violations []model.Violation
	for key, vals := range cb.Required {
		violations = append(violations, FindDuplicates(key, vals)...)
	}
	for key, vals := range cb.Optional {
		violations = append(violations, FindDuplicates(key, vals)...)
	}
	return violations
}
