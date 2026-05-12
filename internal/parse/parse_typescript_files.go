//ff:func feature=parse type=parser control=iteration dimension=1
//ff:what TypeScript 파일 목록을 배치로 파싱하여 TypeScriptFile 슬라이스 반환
package parse

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/park-jun-woo/filefunc/internal/model"
)

//go:embed ts_ast.js
var tsAstScript []byte

// ParseTypeScriptFiles parses multiple TypeScript files in batch mode via ts_ast.js.
func ParseTypeScriptFiles(paths []string, projectRoot string) ([]*model.TypeScriptFile, error) {
	if len(paths) == 0 {
		return nil, nil
	}

	tmpDir, err := os.MkdirTemp("", "filefunc-tsast-*")
	if err != nil {
		return nil, fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	scriptPath := filepath.Join(tmpDir, "ts_ast.js")
	if err := os.WriteFile(scriptPath, tsAstScript, 0644); err != nil {
		return nil, fmt.Errorf("write ts_ast.js: %w", err)
	}

	args := append([]string{scriptPath, "--batch"}, paths...)
	cmd := exec.Command("node", args...)
	cmd.Env = append(os.Environ(), "NODE_PATH="+filepath.Join(projectRoot, "node_modules"))
	out, err := cmd.Output()
	if err != nil {
		return nil, wrapTsAstError(err)
	}

	var results []TsAstResult
	if err := json.Unmarshal(out, &results); err != nil {
		return nil, fmt.Errorf("parse ts_ast.js output: %w", err)
	}

	tsFiles := make([]*model.TypeScriptFile, 0, len(results))
	for _, r := range results {
		if r.Error != "" {
			return nil, fmt.Errorf("ts_ast.js error for %s: %s", r.Path, r.Error)
		}
		tsFiles = append(tsFiles, buildTypeScriptFile(r))
	}

	return tsFiles, nil
}
