//ff:func feature=chain type=parser control=iteration dimension=1
//ff:what 같은 부모에게 호출되는 func 목록 (2촌 형제)
package chain

// FindSiblings returns funcs that share the same caller as the given func.
// Excludes the func itself.
func FindSiblings(g *CallGraph, name string) []string {
	seen := make(map[string]bool)
	seen[name] = true
	var result []string
	for _, parent := range g.Parents[name] {
		AddUnique(g.Children[parent], seen, &result)
	}
	return result
}
