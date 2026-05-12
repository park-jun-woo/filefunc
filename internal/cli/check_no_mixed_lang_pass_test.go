//ff:func feature=validate type=util control=sequence
//ff:what test: TestCheckNoMixedLangPass
package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckNoMixedLangPass(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test"), 0644)
	os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main"), 0644)

	err := CheckNoMixedLang(dir, "go", nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
