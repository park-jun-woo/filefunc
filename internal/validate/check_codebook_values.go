//ff:func feature=validate type=rule control=iteration dimension=1
//ff:what A2: 어노테이션 값이 코드북에 존재하는지 검증
//ff:checked llm=gpt-oss:20b hash=fbc8cbdc
package validate

import (
	"fmt"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// CheckCodebookValues checks A2: annotation values must exist in the codebook.
func CheckCodebookValues(gf *model.GoFile, cb *model.Codebook) []model.Violation {
	if cb == nil || gf.Annotation == nil {
		return nil
	}

	var violations []model.Violation
	for key, val := range gf.Annotation.Func {
		allowed := AllowedValues(cb, key)
		if allowed != nil && !Contains(allowed, val) {
			violations = append(violations, model.Violation{
				File:    gf.Path,
				Rule:    "A2",
				Level:   "ERROR",
				Message: fmt.Sprintf("codebook has no %s=%s", key, val),
			})
		}
	}
	for key, val := range gf.Annotation.Type {
		allowed := AllowedValues(cb, key)
		if allowed != nil && !Contains(allowed, val) {
			violations = append(violations, model.Violation{
				File:    gf.Path,
				Rule:    "A2",
				Level:   "ERROR",
				Message: fmt.Sprintf("codebook has no %s=%s", key, val),
			})
		}
	}
	return violations
}
