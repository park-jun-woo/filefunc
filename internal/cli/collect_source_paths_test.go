//ff:func feature=cli type=util control=sequence
//ff:what test: TestCollectSourcePaths
package cli

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func TestCollectSourcePaths(t *testing.T) {
	files := []model.SourceFile{
		&model.PythonFile{Path: "a.py"},
		&model.PythonFile{Path: "b.py"},
	}
	paths := collectSourcePaths(files)
	if len(paths) != 2 {
		t.Fatalf("expected 2 paths, got %d", len(paths))
	}
	if paths[0] != "a.py" || paths[1] != "b.py" {
		t.Errorf("expected [a.py b.py], got %v", paths)
	}
}
