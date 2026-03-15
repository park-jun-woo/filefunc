//ff:func feature=cli type=command
//ff:what Provider로 what-body 일치 점수를 판정
//ff:calls BuildPrompt, ParseScore
//ff:uses Provider
//ff:checked llm=gpt-oss:20b hash=43ab8822
package llm

// VerifyWhat uses an LLM Provider to score how well what matches body.
// Returns the score (0.0~1.0) and any error.
func VerifyWhat(p Provider, what string, body string) (float64, error) {
	prompt := BuildPrompt(what, body)
	response, err := p.Generate(prompt)
	if err != nil {
		return 0, err
	}
	return ParseScore(response)
}
