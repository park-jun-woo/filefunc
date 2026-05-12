//ff:func feature=validate type=rule control=iteration dimension=1
//ff:what eslintFileResult 슬라이스에서 severity 2 메시지를 N4 violation으로 변환
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// collectESLintViolations converts eslint results into N4 violations for error-severity messages.
func collectESLintViolations(results []eslintFileResult) []model.Violation {
	var violations []model.Violation
	for _, r := range results {
		violations = append(violations, extractFileErrors(r)...)
	}
	return violations
}
