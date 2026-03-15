//ff:type feature=cli type=model
//ff:what LLM API 제공자의 인터페이스 정의
package llm

// Provider is the interface for LLM API providers.
type Provider interface {
	Generate(prompt string) (string, error)
}
