//ff:func feature=context type=command control=sequence
//ff:what LLM으로 codebook feature를 선택하여 반환
package context

import "github.com/park-jun-woo/filefunc/internal/model"

// SelectFeature asks LLM to select relevant features from codebook.
func SelectFeature(prompt string, cb *model.Codebook, generate func(string) (string, error)) ([]string, error) {
	llmPrompt := BuildFeaturePrompt(prompt, cb)
	resp, err := generate(llmPrompt)
	if err != nil {
		return nil, err
	}
	features := ParseFeatures(resp)
	if len(features) == 0 {
		return nil, nil
	}
	return features, nil
}
