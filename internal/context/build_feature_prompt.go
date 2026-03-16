//ff:func feature=context type=formatter control=sequence
//ff:what feature+주석 맵으로 feature 선택용 LLM 프롬프트 생성
package context

import "encoding/json"

// BuildFeaturePrompt creates a prompt for LLM to select relevant features.
func BuildFeaturePrompt(prompt string, codebookRaw string) string {
	features := extractFeatureDescriptions(codebookRaw)
	req := map[string]interface{}{
		"task":     "Select the 1~2 most relevant features for the user question. ONLY use values from the feature list. Pick 1 if obvious, 2 if the task spans multiple areas.",
		"features": features,
		"question": prompt,
		"format":   "JSON array of feature strings, max 2",
	}
	b, _ := json.Marshal(req)
	return string(b)
}
