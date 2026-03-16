//ff:func feature=context type=command control=sequence
//ff:what what 배치 스코어링 — GoFile 목록 대상
package context

import "github.com/park-jun-woo/filefunc/internal/model"

// ScoreWhat runs what-based relevance scoring and filters by rate.
func ScoreWhat(files []*model.GoFile, prompt string, rate float64, generate func(string) (string, error)) ([]*model.GoFile, map[int]float64, int, error) {
	batchPrompt := BuildWhatPrompt(prompt, files)
	resp, err := generate(batchPrompt)
	if err != nil {
		return nil, nil, 0, err
	}
	scores := ParseScores(resp)
	return filterByScore(files, scores, rate)
}
