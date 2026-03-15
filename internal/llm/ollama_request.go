//ff:type feature=cli type=model
//ff:what ollama /api/generate 요청 구조체
package llm

type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}
