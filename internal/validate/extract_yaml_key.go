//ff:func feature=validate type=util control=sequence
//ff:what YAML 라인에서 키를 추출 ("validate: ..." → "validate", "- validate" → "validate")
package validate

import "strings"

// extractYAMLKey extracts the key from a YAML line.
func extractYAMLKey(line string) string {
	key := strings.SplitN(line, ":", 2)[0]
	key = strings.TrimPrefix(key, "- ")
	return strings.TrimSpace(key)
}
