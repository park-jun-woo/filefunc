//ff:func feature=chain type=parser control=iteration
//ff:what 단방향 깊이 탐색의 재귀 헬퍼
package chain

func traverseDepthRecur(g *CallGraph, current string, direction string, depth int, maxDepth int, seen map[string]bool, results *[]ChonResult) {
	if depth >= maxDepth {
		return
	}

	var nexts []string
	if direction == "calls" {
		nexts = g.Children[current]
	} else {
		nexts = g.Parents[current]
	}

	for _, next := range nexts {
		if seen[next] {
			continue
		}
		seen[next] = true
		*results = append(*results, ChonResult{next, depth + 1, direction})
		traverseDepthRecur(g, next, direction, depth+1, maxDepth, seen, results)
	}
}
