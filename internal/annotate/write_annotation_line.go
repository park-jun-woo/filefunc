//ff:func feature=annotate type=formatter
//ff:what 파일의 //ff: 라인을 추가/갱신/제거
//ff:calls InsertAfterAnnotations, ReplaceAnnotationLine
//ff:checked llm=gpt-oss:20b hash=da4749c7
package annotate

import (
	"os"
	"strings"
)

// WriteAnnotationLine updates or inserts a //ff:<key> line in a file.
// If value is empty, removes the existing line.
// Returns true if the file was modified.
func WriteAnnotationLine(path string, key string, value string) (bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}

	lines := strings.Split(string(data), "\n")
	prefix := "//ff:" + key + " "
	newLine := prefix + value
	if value == "" {
		newLine = ""
	}

	var result []string
	modified := false
	found := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		replacement, matched := ReplaceAnnotationLine(trimmed, prefix, key, newLine, value)
		if !matched {
			result = append(result, line)
			continue
		}
		found = true
		if replacement != trimmed {
			modified = true
		}
		if replacement != "" {
			result = append(result, replacement)
		}
	}

	if !found && value != "" {
		result = InsertAfterAnnotations(result, newLine)
		modified = true
	}

	if !modified {
		return false, nil
	}

	return true, os.WriteFile(path, []byte(strings.Join(result, "\n")), 0644)
}
