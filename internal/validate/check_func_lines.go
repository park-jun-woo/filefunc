//ff:func feature=validate type=rule control=iteration dimension=1
//ff:what Q2/Q3 toulmin rule — func 라인 수 위반 시 violation 반환
package validate

import (
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// CheckFuncLines returns (true, []model.Violation) if the file violates Q2 or Q3.
// Q2: func > 1000 lines → ERROR.
// Q3: control=sequence func > 100 lines → ERROR.
func CheckFuncLines(claim any, ground any, backing any) (bool, any) {
	sf := ground.(*ValidateGround).File

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, sf.GetPath(), nil, 0)
	if err != nil {
		return false, nil
	}

	q3Limit, q3Applies := Q3Limit(sf)
	var violations []model.Violation

	for _, decl := range f.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok || fd.Body == nil {
			continue
		}
		violations = checkOneFuncLines(fset, sf.GetPath(), fd, q3Limit, q3Applies, violations)
	}

	if len(violations) > 0 {
		return true, violations
	}
	return false, nil
}
