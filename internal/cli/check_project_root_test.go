//ff:func feature=cli type=util control=sequence
//ff:what test: TestCheckProjectRoot
package cli

import (
	"os"
	"testing"
)

func TestCheckProjectRoot(t *testing.T) {
	if err := CheckProjectRoot("/nonexistent/path"); err == nil {
		t.Error("expected error for nonexistent path")
	}

	tmp := t.TempDir()
	if err := CheckProjectRoot(tmp); err == nil {
		t.Error("expected error for missing go.mod")
	}

	os.WriteFile(tmp+"/go.mod", []byte("module test"), 0644)
	if err := CheckProjectRoot(tmp); err == nil {
		t.Error("expected error for missing codebook.yaml")
	}

	os.WriteFile(tmp+"/codebook.yaml", []byte(""), 0644)
	if err := CheckProjectRoot(tmp); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
