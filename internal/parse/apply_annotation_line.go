//ff:func feature=parse type=parser control=selection
//ff:what //ff: key-value 쌍을 Annotation 구조체에 적용
//ff:checked llm=gpt-oss:20b hash=a752d254
package parse

import "github.com/park-jun-woo/filefunc/internal/model"

// ApplyAnnotationLine applies a parsed //ff: key-value pair to an Annotation.
func ApplyAnnotationLine(ann *model.Annotation, key, value string) {
	switch key {
	case "func":
		ann.Func = ParseFuncPairs(value)
	case "type":
		ann.Type = ParseFuncPairs(value)
	case "what":
		ann.What = value
	case "why":
		ann.Why = value
	case "checked":
		ann.Checked = ParseFuncPairs(value)
	}
}
