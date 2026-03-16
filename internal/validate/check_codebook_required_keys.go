//ff:func feature=validate type=rule control=iteration dimension=1
//ff:what C1: codebook required 섹션에 최소 1개 키가 있고 각 키에 최소 1개 값이 있는지 검증
package validate

import (
	"fmt"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// CheckCodebookRequiredKeys checks C1: required section must exist with at least one key,
// and each key must have at least one value.
func CheckCodebookRequiredKeys(cb *model.Codebook) []model.Violation {
	if len(cb.Required) == 0 {
		return []model.Violation{{
			File:    "codebook.yaml",
			Rule:    "C1",
			Level:   "ERROR",
			Message: "codebook must have a required section with at least one key",
		}}
	}

	var violations []model.Violation
	for key, vals := range cb.Required {
		if len(vals) == 0 {
			violations = append(violations, model.Violation{
				File:    "codebook.yaml",
				Rule:    "C1",
				Level:   "ERROR",
				Message: fmt.Sprintf("required key %q must have at least one value", key),
			})
		}
	}
	return violations
}
