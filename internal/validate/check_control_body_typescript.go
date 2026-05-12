//ff:func feature=validate type=rule control=iteration dimension=1
//ff:what Q4 TypeScript 전용 — TypeScriptFile.Q4Violations에서 PURE 10줄 초과 시 violation 반환
package validate

import (
	"fmt"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// CheckControlBodyTypeScript returns (true, []model.Violation) if a TypeScript control body's
// PURE lines exceed 10. Uses pre-computed Q4Violations from ts_ast.js.
func CheckControlBodyTypeScript(claim any, ground any, backing any) (bool, any) {
	sf := ground.(*ValidateGround).File
	tf, ok := sf.(*model.TypeScriptFile)
	if !ok {
		return false, nil
	}

	var violations []model.Violation
	for _, q := range tf.Q4Violations {
		if q.PureLines > 10 {
			violations = append(violations, model.Violation{
				File:  tf.Path,
				Rule:  "Q4",
				Level: "ERROR",
				Message: fmt.Sprintf("func %s: %s body at line %d has %d PURE lines; extract to sequence func (max 10)",
					q.FuncName, q.StmtType, q.Line, q.PureLines),
			})
		}
	}

	if len(violations) > 0 {
		return true, violations
	}
	return false, nil
}
