//ff:func feature=validate type=rule control=iteration dimension=1
//ff:what Q2/Q3 toulmin rule — func 라인 수 위반 시 violation 반환
package validate

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// RuleQ2Q3 returns (true, []model.Violation) if the file violates Q2 or Q3 (func line limits).
func RuleQ2Q3(claim any, ground any) (bool, any) {
	gf := ground.(*ValidateGround).File

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, gf.Path, nil, 0)
	if err != nil {
		return false, nil
	}

	q3Limit := Q3Limit(gf)
	var violations []model.Violation

	for _, decl := range f.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok || fd.Body == nil {
			continue
		}

		startLine := fset.Position(fd.Pos()).Line
		endLine := fset.Position(fd.End()).Line
		lines := endLine - startLine + 1

		if lines > 1000 {
			violations = append(violations, model.Violation{
				File:    gf.Path,
				Rule:    "Q2",
				Level:   "ERROR",
				Message: fmt.Sprintf("func %s is %d lines; maximum is 1000", fd.Name.Name, lines),
			})
			continue
		}

		if lines <= q3Limit {
			continue
		}

		msg := fmt.Sprintf("func %s is %d lines; recommended maximum is %d", fd.Name.Name, lines, q3Limit)
		if HasBacktick(fd) {
			msg += "\n  hint: backtick string detected — extract to a var-only file (e.g. template_xxx.go) to exempt from Q3"
		}
		violations = append(violations, model.Violation{
			File:    gf.Path,
			Rule:    "Q3",
			Level:   "WARNING",
			Message: msg,
		})
	}

	if len(violations) > 0 {
		return true, violations
	}
	return false, nil
}
