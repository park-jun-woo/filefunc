//ff:func feature=cli type=util control=selection
//ff:what chain func의 target을 qualified name으로 해석 — 동명 함수 시 후보 안내
package cli

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/filefunc/internal/chain"
)

func resolveTarget(g *chain.CallGraph, target string, pkg string) (string, error) {
	if strings.Contains(target, ".") {
		return resolveQualifiedTarget(g, target)
	}
	matches := findMatchingFuncs(g, target)
	switch {
	case len(matches) == 0:
		return "", fmt.Errorf("func %q not found in call graph", target)
	case pkg != "":
		return resolveWithPackage(matches, target, pkg)
	case len(matches) == 1:
		return matches[0], nil
	default:
		return "", fmt.Errorf("multiple funcs named %q found:\n  %s\nUse --package to specify", target, strings.Join(matches, "\n  "))
	}
}
