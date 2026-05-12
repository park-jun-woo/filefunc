//ff:func feature=validate type=rule control=sequence
//ff:what eslint --format json 실행 후 린트 위반에 N4 violation 생성; 미설치 시 ERROR
package validate

import (
	"bytes"
	"encoding/json"
	"os/exec"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// CheckESLint runs "eslint --format json" on the given paths and returns N4 violations
// for each error-severity finding. If eslint is not installed, returns N4 ERROR.
func CheckESLint(paths []string) []model.Violation {
	if _, err := exec.LookPath("eslint"); err != nil {
		return []model.Violation{{
			Rule:    "N4",
			Level:   "ERROR",
			Message: "eslint not installed; run: npm install -D eslint",
		}}
	}

	if len(paths) == 0 {
		return nil
	}

	args := append([]string{"--format", "json"}, paths...)
	cmd := exec.Command("eslint", args...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	_ = cmd.Run()

	var results []eslintFileResult
	if err := json.Unmarshal(stdout.Bytes(), &results); err != nil {
		return nil
	}

	return collectESLintViolations(results)
}
