//ff:type feature=parse type=model
//ff:what 파싱된 //ff: 어노테이션 메타데이터를 담는 구조체
package model

// Annotation holds parsed //ff: metadata from a Go source file.
type Annotation struct {
	Func  map[string]string // key-value pairs from //ff:func (e.g. feature=validate type=rule)
	Type  map[string]string // key-value pairs from //ff:type (e.g. feature=validate type=model)
	What    string            // //ff:what — what this func/type does (required)
	Why     string            // //ff:why — why it was designed this way (optional)
	Checked map[string]string // //ff:checked — LLM verification signature (llm=model hash=xxx)
}
