//ff:func feature=cli type=loader control=sequence
//ff:what ollama에 모델이 존재하는지 확인하고 없으면 pull 여부를 질의
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
	answer, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("read input: %w", err)
	}
	if strings.TrimSpace(strings.ToLower(answer)) != "y" {
		return fmt.Errorf("model %s not available", model)
	}

	return PullModel(endpoint, model)
}
