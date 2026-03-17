//ff:func feature=cli type=util control=iteration dimension=1
//ff:what 그래프에서 함수명이 일치하는 모든 qualified name을 수집
package cli

import (
	"sort"

	"github.com/park-jun-woo/filefunc/internal/chain"
)

func findMatchingFuncs(g *chain.CallGraph, target string) []string {
	seen := make(map[string]bool)
	for key := range g.Children {
		if nameFromQualified(key) == target {
			seen[key] = true
		}
	}
	for key := range g.Parents {
		if nameFromQualified(key) == target {
			seen[key] = true
		}
	}
	var matches []string
	for k := range seen {
		matches = append(matches, k)
	}
	sort.Strings(matches)
	return matches
}
