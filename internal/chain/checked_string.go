//ff:func feature=chain type=formatter control=iteration dimension=1
//ff:what Checked 맵을 "key=value key=value" 형태의 문자열로 변환
package chain

import (
	"fmt"
	"sort"
	"strings"
)

// checkedString converts a checked map to a space-separated key=value string.
func checkedString(m map[string]string) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, m[k]))
	}
	return strings.Join(parts, " ")
}
