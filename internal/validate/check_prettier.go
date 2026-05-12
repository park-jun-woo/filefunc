//ff:func feature=validate type=rule control=iteration dimension=1
//ff:what prettier --check 실행 후 포매팅 위반 파일에 N4 violation 생성; 미설치 시 ERROR
package validate

import (
	"bufio"
	"bytes"
	"os/exec"
	"strings"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// CheckPrettier runs "prettier --check" on the given paths and returns N4 violations
// for files that are not formatted. If prettier is not installed, returns N4 ERROR.
func CheckPrettier(paths []string) []model.Violation {
	if _, err := exec.LookPath("prettier"); err != nil {
		return []model.Violation{{
			Rule:    "N4",
			Level:   "ERROR",
			Message: "prettier not installed; run: npm install -D prettier",
		}}
	}

	if len(paths) == 0 {
		return nil
	}

	args := append([]string{"--check"}, paths...)
	cmd := exec.Command("prettier", args...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err == nil {
		return nil
	}

	var violations []model.Violation
	scanner := bufio.NewScanner(&stdout)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		file := parsePrettierLine(line)
		if file == "" {
			continue
		}
		violations = append(violations, model.Violation{
			File:    file,
			Rule:    "N4",
			Level:   "ERROR",
			Message: "not formatted; run: prettier --write " + file,
		})
	}
	return violations
}
