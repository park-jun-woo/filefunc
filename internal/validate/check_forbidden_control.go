//ff:func feature=validate type=util control=iteration dimension=1
//ff:what 금지된 제어 구조가 depth 1에 존재하는지 검사
package validate

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/filefunc/internal/model"
	"github.com/park-jun-woo/filefunc/internal/parse"
)

func checkForbiddenControl(gf *model.GoFile, b *ControlMatchBacking) (bool, any) {
	for _, f := range strings.Split(b.MustNotHave, "|") {
		if !hasForbidden(gf.Path, f) {
			continue
		}
		msg := b.Message
		if strings.Contains(msg, "%s") {
			actual := parse.DetectControl(gf.Path)
			msg = fmt.Sprintf(b.Message, actual, actual)
		}
		return true, []model.Violation{{
			File:    gf.Path,
			Rule:    b.Rule,
			Level:   "ERROR",
			Message: msg,
		}}
	}
	return false, nil
}
