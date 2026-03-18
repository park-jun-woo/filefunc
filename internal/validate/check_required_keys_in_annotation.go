//ff:func feature=validate type=rule control=iteration dimension=1
//ff:what A8: 어노테이션에 codebook required 키가 모두 존재하는지 검증
package validate

import (
	"fmt"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// CheckRequiredKeysInAnnotation checks A8: all required codebook keys must be
// present in the //ff:func or //ff:type annotation.
func CheckRequiredKeysInAnnotation(gf *model.GoFile, cb *model.Codebook) []model.Violation {
	if cb == nil || gf.Annotation == nil {
		return nil
	}

	meta := gf.Annotation.Func
	if len(meta) == 0 {
		meta = gf.Annotation.Type
	}
	if len(meta) == 0 {
		return nil
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
	return violations
}
