//ff:func feature=validate type=util control=sequence
//ff:what control=sequence인 경우에만 Q3 적용 여부와 줄 수 제한을 반환
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// Q3Limit returns (100, true) if the file is control=sequence, (0, false) otherwise.
// Q3 only applies to sequence funcs. iteration/selection are governed by Q4 (control body limit).
func Q3Limit(gf *model.GoFile) (int, bool) {
	if gf.Annotation != nil && gf.Annotation.Func["control"] == "sequence" {
		return 100, true
	}
	return 0, false
}
