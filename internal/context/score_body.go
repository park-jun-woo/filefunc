//ff:func feature=context type=command control=sequence
//ff:what 본문 배치 스코어링 — GoFile 목록 대상
package context

import "github.com/park-jun-woo/filefunc/internal/model"

// ScoreBody runs body-based relevance scoring and filters by rate.
func ScoreBody(files []*model.GoFile, prompt string, rate float64, generate func(string) (string, error)) ([]*model.GoFile, map[int]float64, int, error) {
	batchPrompt, indices := BuildBodyPrompt(prompt, files)
	if len(indices) == 0 {
		return files, defaultScores(files), 0, nil
	}
	resp, err := generate(batchPrompt)
	if err != nil {
		return nil, nil, 0, err
	}
	scores := ParseScores(resp)
	return filterByScoreWithIndices(files, indices, scores, rate)
}
