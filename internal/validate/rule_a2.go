//ff:func feature=validate type=rule control=iteration dimension=1
//ff:what A2 toulmin rule — 어노테이션 값이 코드북에 존재하지 않으면 violation 반환
package validate

import (
	"fmt"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// RuleA2 returns (true, []model.Violation) if the file violates A2 (annotation value not in codebook).
func RuleA2(claim any, ground any) (bool, any) {
	g := ground.(*ValidateGround)
	gf := g.File
	cb := g.Codebook
	if cb == nil || gf.Annotation == nil {
		return false, nil
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

	if len(violations) > 0 {
		return true, violations
	}
	return false, nil
}
