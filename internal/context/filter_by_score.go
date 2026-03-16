//ff:func feature=context type=util control=iteration dimension=1
//ff:what 인덱스별 점수와 임계값으로 ChonResult를 필터링
package context

import "github.com/park-jun-woo/filefunc/internal/chain"

// filterByScore keeps chon=1 unconditionally and filters chon=2+ by score threshold.
func filterByScore(results []chain.ChonResult, indices []int, scores []float64, rate float64) ([]chain.ChonResult, map[int]float64, int, error) {
	scoreMap := make(map[int]float64)
	for j, idx := range indices {
		if j < len(scores) {
			scoreMap[idx] = scores[j]
		}
	}
	var kept []chain.ChonResult
	keptScores := make(map[int]float64)
	removed := 0
	for i, r := range results {
		if r.Chon <= 1 {
			keptScores[len(kept)] = 1.0
			kept = append(kept, r)
			continue
		}
		s, ok := scoreMap[i]
		if !ok {
			removed++
			continue
		}
		if s >= rate {
			keptScores[len(kept)] = s
			kept = append(kept, r)
		} else {
			removed++
		}
	}
	return kept, keptScores, removed, nil
}
