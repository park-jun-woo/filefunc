//ff:func feature=validate type=util control=iteration dimension=1
//ff:what Q4 switch의 각 case 절 PURE body를 검사하여 위반 시 violation 추가
package validate

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func checkQ4Cases(fset *token.FileSet, path, funcName string, body *ast.BlockStmt, violations []model.Violation) []model.Violation {
	if body == nil {
		return violations
	}
	for _, stmt := range body.List {
		cc, ok := stmt.(*ast.CaseClause)
		if !ok {
			continue
		}
		total := q4CaseBodyLines(fset, cc)
		inner := q4InnerControlLines(fset, cc.Body)
		pure := total - inner
		if pure > 10 {
			line := fset.Position(cc.Pos()).Line
			violations = append(violations, model.Violation{
				File:  path,
				Rule:  "Q4",
				Level: "ERROR",
				Message: fmt.Sprintf(
					"func %s: case at line %d has %d pure lines (total %d - control %d); extract to sequence func (max 10)",
					funcName, line, pure, total, inner),
			})
		}
	}
	return violations
}
