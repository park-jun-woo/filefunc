//ff:func feature=context type=parser control=sequence
//ff:what 한 줄에서 점수를 추출 — "0.85" 또는 "1. 0.85" 형식
package context

import "strconv"

// parseSingleScore extracts a score from a line. Returns -1 if not parseable.
func parseSingleScore(line string) float64 {
	if m := numberedPattern.FindStringSubmatch(line); m != nil {
		s, err := strconv.ParseFloat(m[2], 64)
		if err != nil {
			return -1
		}
		return s
	}
	if m := plainPattern.FindStringSubmatch(line); m != nil {
		s, err := strconv.ParseFloat(m[1], 64)
		if err != nil {
			return -1
		}
		return s
	}
	return -1
}
