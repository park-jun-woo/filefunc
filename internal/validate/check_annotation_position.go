//ff:func feature=validate type=rule control=iteration
//ff:what A6: //ff: 어노테이션이 파일 최상단에 위치하는지 검증
//ff:checked llm=gpt-oss:20b hash=379025dd
package validate

import (
	"bufio"
	"os"
	"strings"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// CheckAnnotationPosition checks A6: //ff: annotations must be at the top of the file,
// before any func/type declarations.
func CheckAnnotationPosition(gf *model.GoFile) []model.Violation {
	f, err := os.Open(gf.Path)
	if err != nil {
		return nil
	}
	defer f.Close()

	seenCode := false
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "//ff:") && seenCode {
			return []model.Violation{{
				File:    gf.Path,
				Rule:    "A6",
				Level:   "ERROR",
				Message: "//ff: annotation must be at the top of the file",
			}}
		}

		if strings.HasPrefix(line, "func ") || strings.HasPrefix(line, "type ") || strings.HasPrefix(line, "var ") {
			seenCode = true
		}
	}
	return nil
}
