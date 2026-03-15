//ff:func feature=cli type=loader
//ff:what ollama 엔드포인트에서 모델을 pull
//ff:checked llm=gpt-oss:20b hash=3297f896
package llm

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// PullModel pulls a model from the ollama endpoint.
func PullModel(endpoint, model string) error {
	fmt.Printf("Pulling %s...\n", model)
	reqBody := fmt.Sprintf(`{"name":"%s","stream":false}`, model)
	resp, err := http.Post(endpoint+"/api/pull", "application/json", strings.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("pull failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("pull failed (%d): %s", resp.StatusCode, string(body))
	}

	fmt.Println("Pull complete.")
	return nil
}
