//ff:func feature=validate type=util control=selection
//ff:what 지정된 종류의 제어 구조가 depth 1에 존재하는지 판별
package validate

import "github.com/park-jun-woo/filefunc/internal/parse"

func hasForbidden(path string, kind string) bool {
	switch kind {
	case "switch":
		return parse.HasSwitchAtDepth1(path)
	case "loop":
		return parse.HasLoopAtDepth1(path)
	}
	return false
}
