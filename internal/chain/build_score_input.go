//ff:func feature=chain type=formatter control=sequence
//ff:what ChonResult의 함수명과 what을 결합하여 reranker document 텍스트를 생성
package chain

import "github.com/park-jun-woo/filefunc/internal/model"

// BuildScoreInput creates a document string for reranker scoring: "FuncName: what text".
func BuildScoreInput(name string, fileMap map[string]model.SourceFile) string {
	sf := fileMap[name]
	display := NameFromQualified(name)
	if sf == nil {
		return display
	}
	ann := sf.GetAnnotation()
	if ann == nil || ann.What == "" {
		return display
	}
	return display + ": " + ann.What
}
