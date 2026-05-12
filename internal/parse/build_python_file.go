//ff:func feature=parse type=parser control=iteration dimension=1
//ff:what PyAstResult를 PythonFile로 변환
package parse

import (
	"path/filepath"
	"strings"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func buildPythonFile(r PyAstResult) *model.PythonFile {
	pf := &model.PythonFile{
		Path:             r.Path,
		Funcs:            r.Functions,
		Classes:          r.Classes,
		Methods:          r.Methods,
		HasInitMethod:    r.HasInitMethod,
		Vars:             r.Vars,
		Lines:            r.Lines,
		MaxDepth:         r.MaxDepth,
		IsTest:           strings.HasPrefix(filepath.Base(r.Path), "test_") || strings.HasSuffix(strings.TrimSuffix(r.Path, ".py"), "_test"),
		Control:          r.Control,
		HasLoopAtDepth1:  r.HasLoopAtDepth1,
		HasMatchAtDepth1: r.HasMatchAtD1,
		FuncLines:        r.FuncLines,
		Calls:            r.Calls,
		BodyHash:         r.BodyHash,
	}

	for _, q := range r.Q4Results {
		pf.Q4Violations = append(pf.Q4Violations, model.Q4Result{
			FuncName:  q.FuncName,
			StmtType:  q.StmtType,
			PureLines: q.PureLines,
			Line:      q.Line,
		})
	}

	return pf
}
