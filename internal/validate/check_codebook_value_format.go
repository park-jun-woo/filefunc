//ff:func feature=validate type=rule
//ff:what C3: codebook 값이 소문자+하이픈 형식인지 검증
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// CheckCodebookValueFormat checks C3: all values must match [a-z][a-z0-9-]*.
func CheckCodebookValueFormat(cb *model.Codebook) []model.Violation {
	var violations []model.Violation
	violations = append(violations, CheckValuesFormat("feature", cb.Feature)...)
	violations = append(violations, CheckValuesFormat("type", cb.Type)...)
	violations = append(violations, CheckValuesFormat("pattern", cb.Pattern)...)
	violations = append(violations, CheckValuesFormat("level", cb.Level)...)
	return violations
}
