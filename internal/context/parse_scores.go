//ff:func feature=context type=parser control=iteration dimension=1
//ff:what LLM 응답에서 "번호. 점수" 형식의 점수 목록을 파싱
package context

import (
	"regexp"
	"strconv"
	"strings"
)

var scorePattern = regexp.MustCompile(`^\d+\.\s*([\d.]+)`)

// ParseScores parses LLM response like "1. 0.8\n2. 0.3\n..." into a score slice.
func ParseScores(response string) []float64 {
	var scores []float64
	for _, line := range strings.Split(strings.TrimSpace(response), "\n") {
		line = strings.TrimSpace(line)
		m := scorePattern.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		s, err := strconv.ParseFloat(m[1], 64)
		if err != nil {
			scores = append(scores, 0)
			continue
		}
		scores = append(scores, s)
	}
	return scores
}
