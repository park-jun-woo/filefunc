//ff:func feature=parse type=parser control=iteration dimension=1
//ff:what TsAstResult를 TypeScriptFile로 변환
package parse

import (
	"path/filepath"
	"strings"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func buildTypeScriptFile(r TsAstResult) *model.TypeScriptFile {
	base := filepath.Base(r.Path)
	isTest := strings.HasSuffix(base, ".test.ts") ||
		strings.HasSuffix(base, ".spec.ts") ||
		strings.HasSuffix(base, ".test.tsx") ||
		strings.HasSuffix(base, ".spec.tsx")

	tf := &model.TypeScriptFile{
		Path:              r.Path,
		Funcs:             r.Functions,
		Classes:           r.Classes,
		Interfaces:        r.Interfaces,
		TypeAliases:       r.TypeAliases,
		Methods:           r.Methods,
		HasConstructor:    r.HasConstructor,
		Vars:              r.Vars,
		Lines:             r.Lines,
		MaxDepth:          r.MaxDepth,
		IsTest:            isTest,
		Control:           r.Control,
		HasLoopAtDepth1:   r.HasLoopAtDepth1,
		HasSwitchAtDepth1: r.HasSwitchAtD1,
		FuncLines:         r.FuncLines,
		Calls:             r.Calls,
		BodyHash:          r.BodyHash,
	}

	for _, q := range r.Q4Results {
		tf.Q4Violations = append(tf.Q4Violations, model.Q4Result{
			FuncName:  q.FuncName,
			StmtType:  q.StmtType,
			PureLines: q.PureLines,
			Line:      q.Line,
		})
	}

	return tf
}
