//ff:func feature=context type=command control=iteration dimension=1
//ff:what what 개별 스코어링 — 각 GoFile의 what을 개별 LLM 호출로 평가
package context

import (
	"fmt"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// ScoreWhat runs what-based relevance scoring per file and filters by rate.
func ScoreWhat(files []*model.GoFile, prompt string, rate float64, generate func(string) (string, error)) ([]*model.GoFile, map[int]float64, int, error) {
	var kept []*model.GoFile
	keptScores := make(map[int]float64)
	removed := 0
	for _, gf := range files {
		name := funcName(gf)
		what := ""
		if gf.Annotation != nil {
			what = gf.Annotation.What
		}
		p := fmt.Sprintf("사용자 작업과 함수의 관련도를 0.0~1.0으로 평가. 직접 수정 대상 0.8+, 영향 0.4~0.7, 무관 0.2 이하.\n\n작업: \"%s\"\n함수: %s: \"%s\"\n\n점수만 출력. 예: 0.80", prompt, name, what)
		resp, err := generate(p)
		if err != nil {
			return nil, nil, 0, err
		}
		scores := ParseScores(resp)
		s := 0.0
		if len(scores) > 0 {
			s = scores[0]
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
