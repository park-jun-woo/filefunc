//ff:func feature=validate type=rule control=iteration dimension=1
//ff:what ruff check 실행 후 린트 위반에 N4 violation 생성; 미설치 시 skip
package validate

import (
	"bufio"
	"bytes"
	"os/exec"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// CheckRuff runs "ruff check --quiet" on the given paths and returns N4 violations
// for each lint finding. If ruff is not installed, returns nil (skip).
func CheckRuff(paths []string) []model.Violation {
	if _, err := exec.LookPath("ruff"); err != nil {
		return nil
	}

	if len(paths) == 0 {
		return nil
	}

	args := append([]string{"check", "--quiet"}, paths...)
	cmd := exec.Command("ruff", args...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err == nil {
		return nil
	}

	var violations []model.Violation
	scanner := bufio.NewScanner(&stdout)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		file, msg := parseRuffLine(line)
		violations = append(violations, model.Violation{
			File:    file,
			Rule:    "N4",
			Level:   "ERROR",
			Message: msg,
		})
	}
	return violations
}
