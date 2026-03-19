//ff:func feature=validate type=util control=iteration dimension=1
//ff:what 코드북 required 키가 어노테이션에 존재하는지 검증
package validate

import (
	"fmt"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func checkCodebookInAnnotation(gf *model.GoFile, cb *model.Codebook, rule string) (bool, any) {
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
				Rule:    rule,
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
