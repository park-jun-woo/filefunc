package walk

import (
	"os"
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

func TestMatchFFIgnore_NoPatterns(t *testing.T) {
	if MatchFFIgnore("src/main.go", "main.go", false, nil) {
		t.Error("expected false with no patterns")
	}
}

func TestMatchFFIgnore_GlobPattern(t *testing.T) {
	patterns := []string{"*.json", "vendor/"}
	if !MatchFFIgnore("config.json", "config.json", false, patterns) {
		t.Error("expected *.json to match config.json")
	}
	if MatchFFIgnore("main.go", "main.go", false, patterns) {
		t.Error("expected *.json not to match main.go")
	}
}

func TestMatchFFIgnore_DirPattern(t *testing.T) {
	patterns := []string{"vendor/"}
	if !MatchFFIgnore("vendor", "vendor", true, patterns) {
		t.Error("expected vendor/ to match vendor dir")
	}
	if MatchFFIgnore("vendor", "vendor", false, patterns) {
		t.Error("expected vendor/ not to match vendor file")
	}
}

func TestMatchPattern_Glob(t *testing.T) {
	if !matchPattern("test.go", "test.go", false, "*.go") {
		t.Error("expected *.go to match test.go")
	}
	if matchPattern("test.txt", "test.txt", false, "*.go") {
		t.Error("expected *.go not to match test.txt")
	}
}

func TestParseFFIgnore_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/.ffignore"
	os.WriteFile(path, []byte("# comment\n\nvendor/\n*.tmp\n"), 0644)
	patterns := ParseFFIgnore(path)
	if len(patterns) != 2 {
		t.Fatalf("len = %d, want 2", len(patterns))
	}
	if patterns[0] != "vendor/" || patterns[1] != "*.tmp" {
		t.Errorf("patterns = %v, want [vendor/ *.tmp]", patterns)
	}
}

func TestParseFFIgnore_NonExistent(t *testing.T) {
	patterns := ParseFFIgnore("/nonexistent/.ffignore")
	if patterns != nil {
		t.Errorf("expected nil, got %v", patterns)
	}
}
