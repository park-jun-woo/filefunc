//ff:func feature=parse type=parser
//ff:what ValueSpec에서 변수명 목록을 추출
//ff:checked llm=gpt-oss:20b hash=21259418
package parse

import "go/ast"

// CollectVarNames extracts variable names from a ValueSpec.
func CollectVarNames(s *ast.ValueSpec) []string {
	var names []string
	for _, name := range s.Names {
		names = append(names, name.Name)
	}
	return names
}
