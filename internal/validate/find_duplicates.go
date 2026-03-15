//ff:func feature=validate type=util
//ff:what 문자열 슬라이스에서 중복 값을 찾아 Violation 목록으로 반환
package validate

import (
	"fmt"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// FindDuplicates returns violations for duplicate values in a string slice.
func FindDuplicates(key string, values []string) []model.Violation {
	seen := make(map[string]bool)
	var violations []model.Violation
	for _, v := range values {
		if seen[v] {
			violations = append(violations, model.Violation{
				File:    "codebook.yaml",
				Rule:    "C2",
				Level:   "ERROR",
				Message: fmt.Sprintf("duplicate value %q in %s", v, key),
			})
		}
		seen[v] = true
	}
	return violations
}
