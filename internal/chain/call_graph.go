//ff:type feature=chain type=model
//ff:what 양방향 호출 관계를 담는 그래프 구조체
package chain

// CallGraph holds bidirectional call relationships.
type CallGraph struct {
	Children map[string][]string // func → callees
	Parents  map[string][]string // func → callers
}
