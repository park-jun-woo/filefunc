//ff:func feature=chain type=parser
//ff:what 촌수 기반 그래프 탐색 — 최대 3촌까지
package chain

// TraverseChon traverses the call graph from a starting func up to maxChon distance.
func TraverseChon(g *CallGraph, start string, maxChon int) []ChonResult {
	var results []ChonResult
	seen := map[string]bool{start: true}

	results = append(results, CollectChon(g.Children[start], 1, "calls", seen)...)
	results = append(results, CollectChon(g.Parents[start], 1, "called-by", seen)...)

	if maxChon < 2 {
		return results
	}

	grandparents := ExpandThrough(g.Parents[start], func(p string) []string { return g.Parents[p] })
	grandchildren := ExpandThrough(g.Children[start], func(c string) []string { return g.Children[c] })
	results = append(results, CollectChon(grandparents, 2, "called-by-2depth", seen)...)
	results = append(results, CollectChon(grandchildren, 2, "calls-2depth", seen)...)
	results = append(results, CollectChon(FindSiblings(g, start), 2, "co-called", seen)...)

	if maxChon < 3 {
		return results
	}

	uncles := ExpandThrough(g.Parents[start], func(p string) []string { return FindSiblings(g, p) })
	nephews := ExpandThrough(FindSiblings(g, start), func(s string) []string { return g.Children[s] })
	results = append(results, CollectChon(uncles, 3, "caller-peer", seen)...)
	results = append(results, CollectChon(nephews, 3, "peer-calls", seen)...)

	return results
}
