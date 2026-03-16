//ff:func feature=chain type=util control=iteration dimension=1
//ff:what ChonResult 목록에서 chon >= 2인 결과 수를 카운트
package chain

// countChon2Plus returns the number of results with chon >= 2.
func countChon2Plus(results []ChonResult) int {
	count := 0
	for _, r := range results {
		if r.Chon >= 2 {
			count++
		}
	}
	return count
}
