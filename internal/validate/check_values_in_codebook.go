//ff:func feature=validate type=util control=iteration dimension=1
//ff:what 어노테이션 값이 코드북에 존재하는지 검증
package validate

import (
	"fmt"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func checkValuesInCodebook(sf model.SourceFile, cb *model.Codebook, rule string) (bool, any) {
	ann := sf.GetAnnotation()
	var violations []model.Violation
	for key, val := range ann.Func {
		allowed := AllowedValues(cb, key)
		if allowed != nil && !Contains(allowed, val) {
			violations = append(violations, model.Violation{
				File:    sf.GetPath(),
				Rule:    rule,
				Level:   "ERROR",
				Message: fmt.Sprintf("codebook has no %s=%s", key, val),
			})
		}
	}
	for key, val := range ann.Type {
		allowed := AllowedValues(cb, key)
		if allowed != nil && !Contains(allowed, val) {
			violations = append(violations, model.Violation{
				File:    sf.GetPath(),
				Rule:    rule,
				Level:   "ERROR",
				Message: fmt.Sprintf("codebook has no %s=%s", key, val),
			})
		}
	}
	if len(violations) > 0 {
		return true, violations
	}
	return false, nil
}
