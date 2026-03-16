//ff:func feature=validate type=util control=iteration dimension=1
//ff:what map의 각 키가 소문자+하이픈 형식인지 검증
package validate

import (
	"fmt"
	"regexp"

	"github.com/park-jun-woo/filefunc/internal/model"
)

var validValuePattern = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)

// CheckValuesFormat returns violations for map keys not matching [a-z][a-z0-9-]*.
func CheckValuesFormat(key string, values map[string]string) []model.Violation {
	var violations []model.Violation
	for v := range values {
		if !validValuePattern.MatchString(v) {
			violations = append(violations, model.Violation{
				File:    "codebook.yaml",
				Rule:    "C3",
				Level:   "ERROR",
				Message: fmt.Sprintf("invalid value %q in %s (must be lowercase + hyphens only)", v, key),
			})
		}
	}
	return violations
}
