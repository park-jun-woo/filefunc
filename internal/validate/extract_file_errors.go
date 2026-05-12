//ff:func feature=validate type=rule control=iteration dimension=1
//ff:what 단일 eslintFileResult에서 severity 2 메시지를 N4 violation으로 추출
package validate

import (
	"fmt"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// extractFileErrors extracts N4 violations from a single eslint file result for error-severity messages.
func extractFileErrors(r eslintFileResult) []model.Violation {
	var violations []model.Violation
	for _, m := range r.Messages {
		if m.Severity != 2 {
			continue
		}
		ruleID := m.RuleID
		if ruleID == "" {
			ruleID = "unknown"
		}
		violations = append(violations, model.Violation{
			File:    r.FilePath,
			Rule:    "N4",
			Level:   "ERROR",
			Message: fmt.Sprintf("eslint %s: %s", ruleID, m.Message),
		})
	}
	return violations
}
