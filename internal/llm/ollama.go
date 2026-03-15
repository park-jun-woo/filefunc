//ff:type feature=cli type=model
//ff:what ollama Provider 구현체
package llm

// OllamaProvider implements Provider using the ollama /api/generate endpoint.
type OllamaProvider struct {
	Endpoint string
	Model    string
}
