//ff:func feature=chain type=parser control=iteration dimension=1
//ff:what --meta 플래그 문자열을 파싱하여 요청된 메타 필드 셋을 반환
package chain

import "strings"

// ParseMetaFlags parses a comma-separated meta flag string into a set.
// "all" expands to meta,what,why,checked. Empty string returns nil.
func ParseMetaFlags(raw string) map[string]bool {
	if raw == "" {
		return nil
	}
	if raw == "all" {
		return map[string]bool{"meta": true, "what": true, "why": true, "checked": true}
	}
	flags := make(map[string]bool)
	for _, f := range strings.Split(raw, ",") {
		flags[strings.TrimSpace(f)] = true
	}
	return flags
}
