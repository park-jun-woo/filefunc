//ff:func feature=validate type=parser control=sequence
//ff:what prettier --check 출력 한 줄에서 포매팅 안 된 파일 경로를 추출; 대상 아니면 빈 문자열 반환
package validate

import "strings"

// parsePrettierLine extracts a file path from a prettier --check output line.
// Returns empty string if the line is not a file violation.
func parsePrettierLine(line string) string {
	if !strings.Contains(line, "[warn]") {
		return ""
	}
	file := strings.TrimSpace(strings.Replace(line, "[warn]", "", 1))
	if file == "" {
		return ""
	}
	if strings.Contains(file, "Code style issues") || strings.Contains(file, "above file") || strings.Contains(file, "Run Prettier") {
		return ""
	}
	return file
}
