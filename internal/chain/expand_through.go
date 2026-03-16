//ff:func feature=chain type=util control=iteration
//ff:what 중간 노드를 거쳐 도달하는 func 목록 수집 (2촌+, 3촌용)
package chain

// ExpandThrough collects funcs reachable by going through intermediates.
func ExpandThrough(intermediates []string, lookup func(string) []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, mid := range intermediates {
		AddUnique(lookup(mid), seen, &result)
	}
	return result
}
