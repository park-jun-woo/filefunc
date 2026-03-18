//ff:func feature=chain type=util control=iteration dimension=1
//ff:what 점수 맵과 임계값으로 ChonResult를 필터링하여 통과 결과 + 제거 건수 반환
package chain

// FilterByRate keeps chon=1 results unconditionally and filters chon=2+ by score threshold.
// Returns filtered results, scores for kept results, and count of removed results.
func FilterByRate(results []ChonResult, scores map[int]float64, rate float64) ([]ChonResult, map[int]float64, int) {
	var kept []ChonResult
	keptScores := make(map[int]float64)
	removed := 0
	for i, r := range results {
		if r.Chon <= 1 {
			kept = append(kept, r)
			continue
		}
		s, ok := scores[i]
		if !ok {
			removed++
			continue
		}
		if s >= rate {
			kept = append(kept, r)
			keptScores[len(kept)-1] = s
		} else {
			removed++
		}
	}
	return kept, keptScores, removed
}
