//ff:func feature=parse type=util control=iteration dimension=1
//ff:what test: TestWalkFilesPy
package walk

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWalkFilesPy(t *testing.T) {
	dir := t.TempDir()
	pyDir := filepath.Join(dir, "src")
	if err := os.MkdirAll(pyDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pyDir, "main.py"), []byte("pass"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pyDir, "util.go"), []byte("package u"), 0644); err != nil {
		t.Fatal(err)
	}

	files, err := WalkFiles(dir, ".py", nil)
	if err != nil {
		t.Fatalf("WalkFiles failed: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}
	for _, f := range files {
		if !strings.HasSuffix(f, ".py") {
			t.Errorf("unexpected file: %s", f)
		}
	}
}
