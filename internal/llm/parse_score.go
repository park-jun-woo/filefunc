//ff:func feature=cli type=parser
//ff:what LLM 응답 문자열에서 0.0~1.0 점수를 추출
//ff:checked llm=gpt-oss:20b hash=465c2cda
package llm

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseScore extracts a float64 score from an LLM response string.
func ParseScore(response string) (float64, error) {
	trimmed := strings.TrimSpace(response)
	score, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse score from response: %q", trimmed)
	}
	if score < 0 || score > 1 {
		return 0, fmt.Errorf("score out of range [0.0, 1.0]: %f", score)
	}
	return score, nil
}
