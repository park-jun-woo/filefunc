//ff:func feature=cli type=loader
//ff:what ollama에 모델이 존재하는지 확인하고 없으면 pull 여부를 질의
//ff:calls ModelExists, PullModel
//ff:checked llm=gpt-oss:20b hash=a281a0aa
package llm

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// CheckModel verifies the model exists on the ollama endpoint.
// If not found, prompts the user to pull it. Returns error if user declines or pull fails.
func CheckModel(endpoint, model string) error {
	if ModelExists(endpoint, model) {
		return nil
	}

	fmt.Printf("Model %s not found. Pull it? [y/N] ", model)
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	if strings.TrimSpace(strings.ToLower(answer)) != "y" {
		return fmt.Errorf("model %s not available", model)
	}

	return PullModel(endpoint, model)
}
