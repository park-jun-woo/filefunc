//ff:func feature=context type=parser control=sequence
//ff:what "validate    # 코드 구조 룰 검증" 형태의 라인에서 이름과 설명을 분리
package context

import "strings"

// parseFeatureLine splits "validate    # description" into name and description.
func parseFeatureLine(line string) (string, string) {
	parts := strings.SplitN(line, "#", 2)
	name := strings.TrimSpace(parts[0])
	desc := ""
	if len(parts) == 2 {
		desc = strings.TrimSpace(parts[1])
	}
	return name, desc
}
