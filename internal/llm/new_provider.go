//ff:func feature=cli type=loader control=selection
//ff:what provider 문자열로 적절한 Provider 구현체를 생성
//ff:checked llm=gpt-oss:20b hash=83753639
package llm

import "fmt"

// NewProvider creates a Provider implementation based on the provider name.
func NewProvider(provider, endpoint, model string) (Provider, error) {
	switch provider {
	case "ollama":
		return &OllamaProvider{Endpoint: endpoint, Model: model}, nil
	}
	return nil, fmt.Errorf("unsupported provider: %s", provider)
}
