//ff:func feature=validate type=parser control=sequence
//ff:what ruff 출력 한 줄을 파싱하여 파일 경로와 메시지를 반환
package validate

import "strings"

// parseRuffLine parses a single ruff output line into (file, message).
// Format: "file:line:col: RULE message" → file, "ruff RULE message".
func parseRuffLine(line string) (string, string) {
	file := line
	if idx := strings.Index(line, ":"); idx >= 0 {
		file = line[:idx]
	}

	msg := "ruff: " + line
	if parts := strings.SplitN(line, " ", 2); len(parts) >= 2 {
		msg = "ruff " + parts[1]
	}
	return file, msg
}
