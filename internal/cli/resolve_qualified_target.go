//ff:func feature=cli type=util control=sequence
//ff:what pkg.FuncName 형식의 qualified target이 그래프에 존재하는지 확인
package cli

import (
	"fmt"

	"github.com/park-jun-woo/filefunc/internal/chain"
)

func resolveQualifiedTarget(g *chain.CallGraph, target string) (string, error) {
	if _, ok := g.Children[target]; ok {
		return target, nil
	}
	if _, ok := g.Parents[target]; ok {
		return target, nil
	}
	return "", fmt.Errorf("func %q not found in call graph", target)
}
