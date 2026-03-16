//ff:func feature=cli type=command control=sequence
//ff:what OllamaProvider의 Generate 메서드 — ollama /api/generate 호출
package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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
