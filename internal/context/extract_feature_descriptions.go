//ff:func feature=context type=parser control=iteration dimension=1
//ff:what codebook.yaml 원본 텍스트에서 feature 이름+주석을 맵으로 추출
package context

import "strings"

// extractFeatureDescriptions parses codebook.yaml raw text to extract
// feature values with their comments as descriptions.
func extractFeatureDescriptions(raw string) map[string]string {
	result := make(map[string]string)
	inFeature := false
	for _, line := range strings.Split(raw, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "feature:" {
			inFeature = true
			continue
		}
		if !inFeature {
			continue
		}
		if !strings.HasPrefix(trimmed, "- ") {
			inFeature = false
			continue
		}
		name, desc := parseFeatureLine(trimmed[2:])
		result[name] = desc
	}
	return result
}
