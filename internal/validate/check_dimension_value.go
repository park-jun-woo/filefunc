//ff:func feature=validate type=rule control=sequence
//ff:what A16: dimension= 값은 양의 정수여야 함
package validate

import (
	"fmt"
	"strconv"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// CheckDimensionValue checks A16: dimension value must be a positive integer.
func CheckDimensionValue(gf *model.GoFile) []model.Violation {
	if gf.IsTest || len(gf.Funcs) == 0 || gf.Annotation == nil {
		return nil
	}
	dim := gf.Annotation.Func["dimension"]
	if dim == "" {
		return nil
	}
	n, err := strconv.Atoi(dim)
	if err != nil || n < 1 {
		return []model.Violation{{
			File:    gf.Path,
			Rule:    "A16",
			Level:   "ERROR",
			Message: fmt.Sprintf("dimension value must be a positive integer, got %q", dim),
		}}
	}
	return nil
}
