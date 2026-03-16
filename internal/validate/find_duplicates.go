//ff:func feature=validate type=util control=iteration dimension=2
//ff:what YAML 원본 텍스트에서 같은 섹션 내 중복 키를 찾아 Violation으로 반환
package validate

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// FindDuplicateKeys scans raw YAML text for duplicate keys at the same indentation level.
func FindDuplicateKeys(raw string) []model.Violation {
	var violations []model.Violation
	section := ""
	seen := make(map[string]map[string]bool)
	for _, line := range strings.Split(raw, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		indent := len(line) - len(strings.TrimLeft(line, " "))
		if indent == 2 && strings.HasSuffix(trimmed, ":") {
			section = strings.TrimSuffix(trimmed, ":")
			seen[section] = make(map[string]bool)
			continue
		}
		if indent != 4 || section == "" {
			continue
		}
		key := extractYAMLKey(trimmed)
		if seen[section][key] {
			violations = append(violations, model.Violation{
				File:    "codebook.yaml",
				Rule:    "C2",
				Level:   "ERROR",
				Message: fmt.Sprintf("duplicate key %q in %s", key, section),
			})
		}
		seen[section][key] = true
	}
	return violations
}
