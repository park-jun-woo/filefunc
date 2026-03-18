//ff:func feature=cli type=loader control=sequence
//ff:what ollama 엔드포인트에서 모델을 pull
package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// PullModel pulls a model from the ollama endpoint.
func PullModel(endpoint, model string) error {
	fmt.Printf("Pulling %s...\n", model)
	reqBody, err := json.Marshal(map[string]interface{}{"name": model, "stream": false})
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}
	resp, err := http.Post(endpoint+"/api/pull", "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("pull failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("pull failed (%d): read body: %w", resp.StatusCode, err)
		}
		return fmt.Errorf("pull failed (%d): %s", resp.StatusCode, string(body))
	}

	fmt.Println("Pull complete.")
	return nil
}
