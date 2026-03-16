//ff:func feature=context type=parser control=iteration dimension=1
//ff:what LLM 응답에서 "번호. 점수" 형식의 점수 목록을 파싱 (thinking 태그 제거)
package context

import (
	"regexp"
	"strconv"
	"strings"
)

var scorePattern = regexp.MustCompile(`^(\d+)\.\s+([\d.]+)\s*$`)
var thinkPattern = regexp.MustCompile(`(?s)<think>.*?</think>`)

// ParseScores parses LLM response like "1. 0.8\n2. 0.3\n..." into a score slice.
// Strips <think>...</think> blocks before parsing.
func ParseScores(response string) []float64 {
	response = thinkPattern.ReplaceAllString(response, "")
	response = strings.TrimSpace(response)
	var scores []float64
	for _, line := range strings.Split(response, "\n") {
		line = strings.TrimSpace(line)
		m := scorePattern.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		s, err := strconv.ParseFloat(m[2], 64)
		if err != nil {
			scores = append(scores, 0)
			continue
		}
		scores = append(scores, s)
	}
	return scores
}
