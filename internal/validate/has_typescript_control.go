//ff:func feature=validate type=util control=selection
//ff:what TypeScript 파일에서 지정된 제어 구조가 depth 1에 존재하는지 판별
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// hasTypeScriptControl checks if the given control kind exists at depth 1 in a TypeScriptFile.
func hasTypeScriptControl(tf *model.TypeScriptFile, kind string) bool {
	switch kind {
	case "switch":
		return tf.HasSwitchAtDepth1
	case "loop":
		return tf.HasLoopAtDepth1
	}
	return false
}
