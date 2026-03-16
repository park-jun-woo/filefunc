//ff:func feature=context type=parser control=iteration dimension=1
//ff:what LLM 응답에서 점수를 파싱 — "0.85" 단독 또는 "번호. 점수" 형식 지원
package context

import (
	"regexp"
	"strings"
)

var numberedPattern = regexp.MustCompile(`^(\d+)\.\s+([\d.]+)\s*$`)
var plainPattern = regexp.MustCompile(`^([\d.]+)\s*$`)
var thinkPattern = regexp.MustCompile(`(?s)<think>.*?</think>`)

// ParseScores parses LLM response into a score slice.
func ParseScores(response string) []float64 {
	response = thinkPattern.ReplaceAllString(response, "")
	response = strings.TrimSpace(response)
	var scores []float64
	for _, line := range strings.Split(response, "\n") {
		line = strings.TrimSpace(line)
		s := parseSingleScore(line)
		if s >= 0 {
			scores = append(scores, s)
		}
	}
	return scores
}
