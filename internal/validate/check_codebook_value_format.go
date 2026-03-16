//ff:func feature=validate type=rule control=iteration
//ff:what C3: codebook required+optional 전체 값이 소문자+하이픈 형식인지 검증
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// CheckCodebookValueFormat checks C3: all values must match [a-z][a-z0-9-]*.
func CheckCodebookValueFormat(cb *model.Codebook) []model.Violation {
	var violations []model.Violation
	for key, vals := range cb.Required {
		violations = append(violations, CheckValuesFormat(key, vals)...)
	}
	for key, vals := range cb.Optional {
		violations = append(violations, CheckValuesFormat(key, vals)...)
	}
	return violations
}
