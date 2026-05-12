//ff:func feature=validate type=rule control=iteration dimension=1
//ff:what black --check 실행 후 포매팅 위반 파일에 N4 violation 생성; 미설치 시 skip
package validate

import (
	"bufio"
	"bytes"
	"os/exec"
	"strings"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// CheckBlack runs "black --check --quiet" on the given paths and returns N4 violations
// for files that would be reformatted. If black is not installed, returns nil (skip).
func CheckBlack(paths []string) []model.Violation {
	if _, err := exec.LookPath("black"); err != nil {
		return nil
	}

	if len(paths) == 0 {
		return nil
	}

	args := append([]string{"--check", "--quiet"}, paths...)
	cmd := exec.Command("black", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err == nil {
		return nil
	}

	var violations []model.Violation
	scanner := bufio.NewScanner(&stderr)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "would reformat ") {
			file := strings.TrimPrefix(line, "would reformat ")
			violations = append(violations, model.Violation{
				File:    file,
				Rule:    "N4",
				Level:   "ERROR",
				Message: "not formatted; run: black " + file,
			})
		}
	}
	return violations
}
