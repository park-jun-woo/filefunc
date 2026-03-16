//ff:func feature=context type=command control=sequence
//ff:what 본문 배치 스코어링 실행 — 프롬프트 생성 + LLM 호출 + 점수 파싱 + 필터링
package context

import (
	"github.com/park-jun-woo/filefunc/internal/chain"
	"github.com/park-jun-woo/filefunc/internal/model"
)

// ScoreBody runs body-based relevance scoring and filters by rate.
// chon=1 results are always kept with score 1.0.
func ScoreBody(results []chain.ChonResult, prompt string, rate float64, fileMap map[string]*model.GoFile, generate func(string) (string, error)) ([]chain.ChonResult, map[int]float64, int, error) {
	batchPrompt, indices := BuildBodyPrompt(prompt, results, fileMap)
	if len(indices) == 0 {
		return results, chon1Scores(results), 0, nil
	}
	resp, err := generate(batchPrompt)
	if err != nil {
		return nil, nil, 0, err
	}
	scores := ParseScores(resp)
	return filterByScore(results, indices, scores, rate)
}
