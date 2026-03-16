//ff:func feature=validate type=util control=sequence
//ff:what control 값에 따라 Q3 줄 수 제한을 반환
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// Q3Limit returns the Q3 line limit based on the control annotation.
// selection: 300, sequence/iteration/default: 100.
func Q3Limit(gf *model.GoFile) int {
	if gf.Annotation != nil && gf.Annotation.Func["control"] == "selection" {
		return 300
	}
	return 100
}
