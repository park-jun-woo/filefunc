//ff:func feature=validate type=rule control=iteration dimension=1
//ff:what A8 toulmin rule — codebook required 키가 어노테이션에 없으면 violation 반환
package validate

import (
	"fmt"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// RuleA8 returns (true, []model.Violation) if the file violates A8 (missing required keys).
func RuleA8(claim any, ground any) (bool, any) {
	g := ground.(*ValidateGround)
	gf := g.File
	cb := g.Codebook
	if cb == nil || gf.Annotation == nil {
		return false, nil
	}

	meta := gf.Annotation.Func
	if len(meta) == 0 {
		meta = gf.Annotation.Type
	}
	if len(meta) == 0 {
		return false, nil
	}

	var violations []model.Violation
	for key := range cb.Required {
		if _, ok := meta[key]; !ok {
			violations = append(violations, model.Violation{
				File:    gf.Path,
				Rule:    "A8",
				Level:   "ERROR",
				Message: fmt.Sprintf("required codebook key %q missing in annotation", key),
			})
		}
	}

	if len(violations) > 0 {
		return true, violations
	}
	return false, nil
}
