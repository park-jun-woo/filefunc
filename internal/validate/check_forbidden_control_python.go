//ff:func feature=validate type=util control=iteration dimension=1
//ff:what Python 파일에서 금지된 제어 구조가 depth 1에 존재하는지 검사
package validate

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func checkForbiddenControlPython(pf *model.PythonFile, b *ControlMatchBacking) (bool, any) {
	for _, kind := range strings.Split(b.MustNotHave, "|") {
		if !hasPythonControl(pf, kind) {
			continue
		}
		msg := b.Message
		if strings.Contains(msg, "%s") {
			msg = fmt.Sprintf(b.Message, pf.Control, pf.Control)
		}
		return true, []model.Violation{{
			File:    pf.GetPath(),
			Rule:    b.Rule,
			Level:   "ERROR",
			Message: msg,
		}}
	}
	return false, nil
}
