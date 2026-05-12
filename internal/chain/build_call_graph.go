//ff:func feature=chain type=parser control=iteration dimension=1
//ff:what 프로젝트 전체 호출 그래프를 양방향으로 구성
package chain

import (
	"fmt"
	"os"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// BuildCallGraph builds a bidirectional call graph from parsed Go files.
func BuildCallGraph(files []model.SourceFile, modulePath string, projFuncs map[string]string) *CallGraph {
	g := &CallGraph{
		Children: make(map[string][]string),
		Parents:  make(map[string][]string),
	}

	for _, sf := range files {
		if sf.IsTestFile() {
			continue
		}
		name := funcName(sf)
		if name == "" {
			continue
		}
		caller := qualifiedName(sf.GetPackage(), name)
		calls, err := extractCallsForFile(sf, modulePath, projFuncs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[WARN] %s: extract calls failed: %v\n", sf.GetPath(), err)
			continue
		}
		g.Children[caller] = calls
		for _, callee := range calls {
			g.Parents[callee] = append(g.Parents[callee], caller)
		}
	}

	return g
}
