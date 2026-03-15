//ff:type feature=cli type=model
//ff:what ollama /api/tags 응답 구조체
package llm

type ollamaTagsResponse struct {
	Models []struct {
		Name string `json:"name"`
	} `json:"models"`
}
