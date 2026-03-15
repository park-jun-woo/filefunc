//ff:type feature=cli type=model
//ff:what ollama /api/generate 응답 구조체
package llm

type ollamaResponse struct {
	Response string `json:"response"`
}
