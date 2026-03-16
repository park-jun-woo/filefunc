//ff:func feature=context type=parser control=sequence
//ff:what LLM 응답에서 feature JSON 배열을 파싱
package context

import (
	"encoding/json"
	"strings"
)

// ParseFeatures parses LLM response containing a JSON array of feature strings.
func ParseFeatures(response string) []string {
	response = strings.TrimSpace(response)
	if idx := strings.Index(response, "["); idx >= 0 {
		response = response[idx:]
	}
	if idx := strings.LastIndex(response, "]"); idx >= 0 {
		response = response[:idx+1]
	}
	var features []string
	if err := json.Unmarshal([]byte(response), &features); err != nil {
		return nil
	}
	return features
}
