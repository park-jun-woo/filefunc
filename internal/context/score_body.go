//ff:func feature=context type=command control=iteration dimension=1
//ff:what 본문 개별 스코어링 — 각 GoFile의 본문을 개별 LLM 호출로 평가
package context

import (
	"fmt"
	"os"

	"github.com/park-jun-woo/filefunc/internal/model"
	"github.com/park-jun-woo/filefunc/internal/parse"
)

// ScoreBody runs body-based relevance scoring per file and filters by rate.
func ScoreBody(files []*model.GoFile, prompt string, rate float64, generate func(string) (string, error)) ([]*model.GoFile, map[int]float64, int, error) {
	var kept []*model.GoFile
	keptScores := make(map[int]float64)
	removed := 0
	for _, gf := range files {
		src, _ := os.ReadFile(gf.Path)
		body := parse.ExtractFuncSource(gf.Path, src)
		if body == "" {
			removed++
			continue
		}
		p := fmt.Sprintf("사용자 작업과 함수의 관련도를 0.0~1.0으로 평가. 직접 수정 대상 0.8+, 영향 0.4~0.7, 무관 0.2 이하.\n\n작업: \"%s\"\n함수:\n%s\n\n점수만 출력. 예: 0.80", prompt, body)
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
