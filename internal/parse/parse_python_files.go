//ff:func feature=parse type=parser control=iteration dimension=1
//ff:what Python 파일 목록을 배치로 파싱하여 PythonFile 슬라이스 반환
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

//go:embed py_ast.py
var pyAstScript []byte

// ParsePythonFiles parses multiple Python files in batch mode via py_ast.py.
func ParsePythonFiles(paths []string) ([]*model.PythonFile, error) {
	if len(paths) == 0 {
		return nil, nil
	}

	tmpDir, err := os.MkdirTemp("", "filefunc-pyast-*")
	if err != nil {
		return nil, fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	scriptPath := filepath.Join(tmpDir, "py_ast.py")
	if err := os.WriteFile(scriptPath, pyAstScript, 0644); err != nil {
		return nil, fmt.Errorf("write py_ast.py: %w", err)
	}

	args := append([]string{scriptPath, "--batch"}, paths...)
	cmd := exec.Command("python3", args...)
	out, err := cmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("py_ast.py failed: %s", string(ee.Stderr))
		}
		return nil, fmt.Errorf("py_ast.py exec: %w", err)
	}

	var results []PyAstResult
	if err := json.Unmarshal(out, &results); err != nil {
		return nil, fmt.Errorf("parse py_ast.py output: %w", err)
	}

	pyFiles := make([]*model.PythonFile, 0, len(results))
	for _, r := range results {
		if r.Error != "" {
			return nil, fmt.Errorf("py_ast.py error for %s: %s", r.Path, r.Error)
		}
		pyFiles = append(pyFiles, buildPythonFile(r))
	}

	return pyFiles, nil
}
