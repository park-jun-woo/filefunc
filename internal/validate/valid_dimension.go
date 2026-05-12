//ff:func feature=validate type=rule control=sequence
//ff:what A16 toulmin rule — dimension= 값이 양의 정수가 아니면 violation 반환
package validate

import (
	"fmt"
	"strconv"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// RuleA16 returns (true, []model.Violation) if the file violates A16.
func ValidDimension(claim any, ground any, backing any) (bool, any) {
	sf := ground.(*ValidateGround).File
	ann := sf.GetAnnotation()
	if len(sf.GetFuncs()) == 0 || ann == nil {
		return false, nil
	}
	dim := ann.Func["dimension"]
	if dim == "" {
		return false, nil
	}
	n, err := strconv.Atoi(dim)
	if err != nil || n < 1 {
		return true, []model.Violation{{
			File:    sf.GetPath(),
			Rule:    "A16",
			Level:   "ERROR",
			Message: fmt.Sprintf("dimension value must be a positive integer, got %q", dim),
		}}
	}
	return false, nil
}
