//ff:func feature=context type=util control=iteration dimension=1
//ff:what GoFile 목록에 기본 점수 0.0을 부여하는 점수 맵 생성
package context

import "github.com/park-jun-woo/filefunc/internal/model"

// defaultScores returns a score map with 0.0 for all files.
func defaultScores(files []*model.GoFile) map[int]float64 {
	scores := make(map[int]float64)
	for i := range files {
		scores[i] = 0.0
	}
	return scores
}
