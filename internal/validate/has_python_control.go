//ff:func feature=validate type=util control=selection
//ff:what Python 파일에서 지정된 제어 구조가 depth 1에 존재하는지 판별
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// hasPythonControl checks if the given control kind exists at depth 1 in a PythonFile.
func hasPythonControl(pf *model.PythonFile, kind string) bool {
	switch kind {
	case "match":
		return pf.HasMatchAtDepth1
	case "loop":
		return pf.HasLoopAtDepth1
	}
	return false
}
