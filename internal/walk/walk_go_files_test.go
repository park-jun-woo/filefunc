//ff:func feature=parse type=util control=iteration dimension=1
//ff:what test: TestWalkGoFiles
package walk

import (
	"strings"
	"testing"
)

func TestWalkGoFiles(t *testing.T) {
	files, err := WalkGoFiles("../parse", nil)
	if err != nil {
		t.Fatalf("WalkGoFiles failed: %v", err)
	}

	if len(files) == 0 {
		t.Fatal("WalkGoFiles returned no files")
	}

	hasGoFile := false
	for _, f := range files {
		if strings.HasSuffix(f, ".go") {
			hasGoFile = true
		}
	}
	if !hasGoFile {
		t.Error("no .go files found")
	}
}
