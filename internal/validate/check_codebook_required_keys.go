//ff:func feature=validate type=rule
//ff:what C1: codebook에 feature, type 키가 최소 1개 값과 함께 존재하는지 검증
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// CheckCodebookRequiredKeys checks C1: feature and type must have at least one value.
func CheckCodebookRequiredKeys(cb *model.Codebook) []model.Violation {
	var violations []model.Violation
	if len(cb.Feature) == 0 {
		violations = append(violations, model.Violation{
			File:    "codebook.yaml",
			Rule:    "C1",
			Level:   "ERROR",
			Message: "codebook must have at least one feature value",
		})
	}
	if len(cb.Type) == 0 {
		violations = append(violations, model.Violation{
			File:    "codebook.yaml",
			Rule:    "C1",
			Level:   "ERROR",
			Message: "codebook must have at least one type value",
		})
	}
	return violations
}
