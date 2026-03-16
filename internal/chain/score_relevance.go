//ff:func feature=chain type=loader control=iteration dimension=1
//ff:what chon=2+ 결과를 vLLM reranker로 스코어링하여 인덱스별 점수 맵 반환
package chain

import (
	"fmt"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// ScoreRelevance scores each chon=2+ result against the prompt using vLLM reranker.
// Returns a map of result index → score.
func ScoreRelevance(results []ChonResult, prompt string, endpoint string, modelName string, fileMap map[string]*model.GoFile) (map[int]float64, error) {
	scores := make(map[int]float64)
	for i, r := range results {
		if r.Chon <= 1 {
			continue
		}
		doc := BuildScoreInput(r.Name, fileMap)
		score, err := callVLLMScore(endpoint, modelName, prompt, doc)
		if err != nil {
			return nil, fmt.Errorf("score failed for %s: %w", r.Name, err)
		}
		scores[i] = score
	}
	return scores, nil
}
