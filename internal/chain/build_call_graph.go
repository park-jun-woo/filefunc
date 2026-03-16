//ff:func feature=chain type=parser control=iteration
//ff:what 프로젝트 전체 호출 그래프를 양방향으로 구성
package chain

import (
	"github.com/park-jun-woo/filefunc/internal/model"
	"github.com/park-jun-woo/filefunc/internal/parse"
)

// BuildCallGraph builds a bidirectional call graph from parsed Go files.
func BuildCallGraph(files []*model.GoFile, modulePath string, projFuncs map[string]bool) *CallGraph {
	g := &CallGraph{
		Children: make(map[string][]string),
		Parents:  make(map[string][]string),
	}

	for _, gf := range files {
		if gf.IsTest {
			continue
		}
		caller := funcName(gf)
		if caller == "" {
			continue
		}
		calls, err := parse.ExtractCalls(gf.Path, modulePath, projFuncs)
		if err != nil {
			continue
		}
		g.Children[caller] = calls
		for _, callee := range calls {
			g.Parents[callee] = append(g.Parents[callee], caller)
		}
	}

	return g
}
