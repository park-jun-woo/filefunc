//ff:func feature=validate type=command control=sequence
//ff:what codebook.yaml 형식 검증 오케스트레이터 (C1~C4 실행)
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// ValidateCodebook runs all codebook validation rules (C1~C4).
func ValidateCodebook(cb *model.Codebook, codebookPath string) []model.Violation {
	var violations []model.Violation
	violations = append(violations, CheckCodebookRequiredKeys(cb)...)
	violations = append(violations, CheckCodebookDuplicates(codebookPath)...)
	violations = append(violations, CheckCodebookValueFormat(cb)...)
	violations = append(violations, CheckCodebookDescription(cb)...)
	return violations
}
