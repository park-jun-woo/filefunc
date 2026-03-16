//ff:func feature=validate type=rule control=iteration dimension=2
//ff:what C4: required 키의 각 값에 description이 비어있지 않아야 함
package validate

import (
	"fmt"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// CheckCodebookDescription checks C4: required values should have non-empty descriptions.
func CheckCodebookDescription(cb *model.Codebook) []model.Violation {
	var violations []model.Violation
	for key, vals := range cb.Required {
		for name, desc := range vals {
			if desc == "" {
				violations = append(violations, model.Violation{
					File:    "codebook.yaml",
					Rule:    "C4",
					Level:   "WARNING",
					Message: fmt.Sprintf("required value %q in %s has no description", name, key),
				})
			}
		}
	}
	return violations
}
