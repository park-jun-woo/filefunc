//ff:func feature=context type=command control=sequence
//ff:what LLM으로 codebook feature를 선택하여 반환
package context

// SelectFeature asks LLM to select relevant features from codebook.
// codebookRaw is the raw codebook.yaml text (includes comments).
func SelectFeature(prompt string, codebookRaw string, generate func(string) (string, error)) ([]string, error) {
	llmPrompt := BuildFeaturePrompt(prompt, codebookRaw)
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
