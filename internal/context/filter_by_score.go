//ff:func feature=context type=util control=iteration dimension=1
//ff:what 점수와 임계값으로 GoFile 목록을 필터링 (1:1 대응)
package context

import "github.com/park-jun-woo/filefunc/internal/model"

// filterByScore filters GoFiles by score threshold. Scores correspond 1:1 with files.
func filterByScore(files []*model.GoFile, scores []float64, rate float64) ([]*model.GoFile, map[int]float64, int, error) {
	var kept []*model.GoFile
	keptScores := make(map[int]float64)
	removed := 0
	for i, gf := range files {
		s := 0.0
		if i < len(scores) {
			s = scores[i]
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
