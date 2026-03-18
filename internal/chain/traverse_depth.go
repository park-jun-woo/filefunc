//ff:func feature=chain type=parser control=sequence
//ff:what 단방향 깊이 탐색 (child-depth 또는 parent-depth)
package chain

// TraverseDepth traverses the call graph in one direction up to maxDepth.
// direction is "calls" or "called-by".
func TraverseDepth(g *CallGraph, start string, direction string, maxDepth int) []ChonResult {
	var results []ChonResult
	seen := map[string]bool{start: true}
	traverseDepthRecur(g, start, direction, 0, maxDepth, seen, &results)
	return results
}
