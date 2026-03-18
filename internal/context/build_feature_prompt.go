//ff:func feature=context type=formatter control=sequence
//ff:what Codebook의 feature description으로 feature 선택용 LLM 프롬프트 생성
package context

import (
	"encoding/json"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// BuildFeaturePrompt creates a prompt for LLM to select relevant features.
func BuildFeaturePrompt(prompt string, cb *model.Codebook) string {
	features := cb.Required["feature"]
	req := map[string]interface{}{
		"task":     "Select the 1~2 most relevant features for the user question. ONLY use values from the feature list. Pick 1 if obvious, 2 if the task spans multiple areas.",
		"features": features,
		"question": prompt,
		"format":   "JSON array of feature strings, max 2",
	}
	b, err := json.Marshal(req)
	if err != nil {
		return ""
	}
	return string(b)
}
