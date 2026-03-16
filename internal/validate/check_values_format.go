//ff:func feature=validate type=util control=iteration
//ff:what 문자열 슬라이스의 각 값이 소문자+하이픈 형식인지 검증
package validate

import (
	"fmt"
	"regexp"

	"github.com/park-jun-woo/filefunc/internal/model"
)

var validValuePattern = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)

// CheckValuesFormat returns violations for values not matching [a-z][a-z0-9-]*.
func CheckValuesFormat(key string, values []string) []model.Violation {
	var violations []model.Violation
	for _, v := range values {
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
