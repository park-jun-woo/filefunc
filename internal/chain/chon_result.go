//ff:type feature=chain type=model
//ff:what 촌수 탐색 결과를 담는 구조체
package chain

// ChonResult holds the traversal result with chon distance.
type ChonResult struct {
	Name string
	Chon int
	Rel  string // "child", "parent", "sibling", "grandchild", "grandparent", "uncle", "nephew"
}
