//ff:func feature=context type=util control=iteration dimension=1
//ff:what chon=1 결과에 점수 1.0을 부여하는 기본 점수 맵 생성
package context

import "github.com/park-jun-woo/filefunc/internal/chain"

// chon1Scores creates a score map with 1.0 for chon=1 results.
func chon1Scores(results []chain.ChonResult) map[int]float64 {
	scores := make(map[int]float64)
	for i, r := range results {
		if r.Chon <= 1 {
			scores[i] = 1.0
		}
	}
	return scores
}
