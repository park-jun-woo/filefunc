//ff:func feature=chain type=parser control=sequence
//ff:what 특정 func을 호출하는 func 목록을 역방향 맵에서 조회
package chain

// FindCallers returns the list of funcs that call the given func.
func FindCallers(g *CallGraph, name string) []string {
	return g.Parents[name]
}
