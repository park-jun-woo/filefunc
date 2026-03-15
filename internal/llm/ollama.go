//ff:type feature=cli type=model
//ff:what ollama /api/generate를 호출하는 Provider 구현
package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// OllamaProvider implements Provider using the ollama /api/generate endpoint.
type OllamaProvider struct {
	Endpoint string
	Model    string
}

type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type ollamaResponse struct {
	Response string `json:"response"`
}

func (o *OllamaProvider) Generate(prompt string) (string, error) {
	reqBody, err := json.Marshal(ollamaRequest{Model: o.Model, Prompt: prompt, Stream: false})
	if err != nil {
		return "", err
	}

	resp, err := http.Post(o.Endpoint+"/api/generate", "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return "", fmt.Errorf("ollama request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama returned %d: %s", resp.StatusCode, string(body))
	}

	var result ollamaResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	return result.Response, nil
}
