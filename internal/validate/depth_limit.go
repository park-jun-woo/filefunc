//ff:func feature=validate type=util control=sequence
//ff:what control과 dimension으로 Q1 depth 상한을 계산
package validate

import (
	"strconv"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// depthLimit returns the maximum allowed nesting depth for a SourceFile.
// sequence=2, selection=2, iteration=dimension+1.
func depthLimit(sf model.SourceFile) int {
	ann := sf.GetAnnotation()
	if ann == nil {
		return 2
	}
	control := ann.Func["control"]
	if control == "iteration" {
		dim := ann.Func["dimension"]
		n, err := strconv.Atoi(dim)
		if err != nil || n < 1 {
			return 2
		}
		return n + 1
	}
	return 2
}
