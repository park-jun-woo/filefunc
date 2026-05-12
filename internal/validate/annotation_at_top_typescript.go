//ff:func feature=validate type=rule control=iteration dimension=1
//ff:what A6 TypeScript 전용 — 어노테이션이 파일 최상단에 위치하지 않으면 violation 반환
package validate

import (
	"bufio"
	"os"
	"strings"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// AnnotationAtTopTypeScript returns (true, []model.Violation) if a TypeScript file
// has //ff: annotations after code has started.
func AnnotationAtTopTypeScript(claim any, ground any, backing any) (bool, any) {
	sf := ground.(*ValidateGround).File

	f, err := os.Open(sf.GetPath())
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
				File:    sf.GetPath(),
				Rule:    "A6",
				Level:   "ERROR",
				Message: "//ff: annotation must be at the top of the file",
			}}
		}

		if isTypeScriptPreamble(line) {
			continue
		}

		if !seenCode && isTypeScriptCodeStart(line) {
			seenCode = true
		}
	}
	if err := scanner.Err(); err != nil {
		return false, err
	}
	return false, nil
}
