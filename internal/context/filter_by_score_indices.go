//ff:func feature=context type=util control=iteration dimension=1
//ff:what 인덱스 매핑된 점수와 임계값으로 GoFile 목록을 필터링
package context

import "github.com/park-jun-woo/filefunc/internal/model"

// filterByScoreWithIndices filters when scores only cover a subset of files (by indices).
func filterByScoreWithIndices(files []*model.GoFile, indices []int, scores []float64, rate float64) ([]*model.GoFile, map[int]float64, int, error) {
	scoreMap := make(map[int]float64)
	for j, idx := range indices {
		if j < len(scores) {
			scoreMap[idx] = scores[j]
		}
	}
	var kept []*model.GoFile
	keptScores := make(map[int]float64)
	removed := 0
	for i, gf := range files {
		s, ok := scoreMap[i]
		if !ok {
			removed++
			continue
		}
		if s >= rate {
			keptScores[len(kept)] = s
			kept = append(kept, gf)
		} else {
			removed++
		}
	}
	return kept, keptScores, removed, nil
}
