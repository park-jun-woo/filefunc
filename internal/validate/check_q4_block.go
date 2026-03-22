//ff:func feature=validate type=util control=sequence
//ff:what Q4 단일 제어문 블록의 PURE body 줄수를 검사하여 위반 시 violation 추가
package validate

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func checkQ4Block(fset *token.FileSet, path, funcName, kind string, body *ast.BlockStmt, violations []model.Violation) []model.Violation {
	if body == nil {
		return violations
	}
	total := q4BlockBodyLines(fset, body)
	inner := q4InnerControlLines(fset, body.List)
	pure := total - inner
	if pure > 10 {
		line := fset.Position(body.Pos()).Line
		violations = append(violations, model.Violation{
			File:  path,
			Rule:  "Q4",
			Level: "ERROR",
			Message: fmt.Sprintf(
				"func %s: %s body at line %d has %d pure lines (total %d - control %d); extract to sequence func (max 10)",
				funcName, kind, line, pure, total, inner),
		})
	}
	return violations
}
