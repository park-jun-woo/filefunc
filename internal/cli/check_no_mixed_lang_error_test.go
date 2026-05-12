//ff:func feature=validate type=util control=sequence
//ff:what test: TestCheckNoMixedLangError
package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckNoMixedLangError(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test"), 0644)
	os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main"), 0644)
	os.WriteFile(filepath.Join(dir, "script.py"), []byte("print('hello')"), 0644)

	err := CheckNoMixedLang(dir, "go", nil)
	if err == nil {
		t.Fatal("expected error for mixed language project")
	}
	if !strings.Contains(err.Error(), ".py") {
		t.Errorf("error should mention .py, got: %v", err)
	}
}
