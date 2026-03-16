//ff:func feature=chain type=formatter control=iteration dimension=1
//ff:what Annotation의 Func 또는 Type 맵에서 key=value 쌍 목록을 생성
package chain

import (
	"fmt"
	"sort"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// metaPairs returns key=value strings from the Func or Type annotation map.
func metaPairs(ann *model.Annotation) []string {
	m := ann.Func
	if len(m) == 0 {
		m = ann.Type
	}
	if len(m) == 0 {
		return nil
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	pairs := make([]string, 0, len(keys))
	for _, k := range keys {
		pairs = append(pairs, fmt.Sprintf("%s=%s", k, m[k]))
	}
	return pairs
}
