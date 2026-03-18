//ff:func feature=validate type=rule control=iteration dimension=1
//ff:what A6 toulmin rule — 어노테이션이 파일 최상단에 위치하지 않으면 violation 반환
package validate

import (
	"bufio"
	"os"
	"strings"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// RuleA6 returns (true, []model.Violation) if the file violates A6 (annotation not at top).
func RuleA6(claim any, ground any) (bool, any) {
	gf := ground.(*ValidateGround).File

	f, err := os.Open(gf.Path)
	if err != nil {
		return false, nil
	}
	defer f.Close()

	seenCode := false
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "//ff:") && seenCode {
			return true, []model.Violation{{
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
	if err := scanner.Err(); err != nil {
		return false, err
	}
	return false, nil
}
