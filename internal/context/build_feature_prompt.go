//ff:func feature=context type=formatter control=sequence
//ff:what codebook feature 선택용 LLM 프롬프트를 생성
package context

import (
	"encoding/json"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// BuildFeaturePrompt creates a prompt for LLM to select relevant features from codebook.
func BuildFeaturePrompt(prompt string, cb *model.Codebook) string {
	features := cb.Required["feature"]
	req := map[string]interface{}{
		"task":     "Given codebook, select the most relevant features for the user's question. Use only codebook feature values.",
		"codebook": map[string]interface{}{"feature": features},
		"question": prompt,
		"format":   "JSON array of feature strings, max 3",
	}
	b, _ := json.Marshal(req)
	return string(b)
}
