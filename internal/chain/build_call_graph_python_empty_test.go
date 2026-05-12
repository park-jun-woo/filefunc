//ff:func feature=chain type=util control=sequence
//ff:what test: TestBuildCallGraphPythonEmpty
package chain

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func TestBuildCallGraphPythonEmpty(t *testing.T) {
	files := []model.SourceFile{
		&model.PythonFile{
			Module: "myapp.main",
			Path:   "myapp/main.py",
			Funcs:  []string{"main"},
			Calls:  nil,
		},
		&model.PythonFile{
			Module: "myapp.config",
			Path:   "myapp/config.py",
			Funcs:  []string{"load_config"},
			Calls:  []string{},
		},
	}

	g := BuildCallGraph(files, "", nil)

	if len(g.Children["myapp.main.main"]) != 0 {
		t.Errorf("main children = %v, want empty", g.Children["myapp.main.main"])
	}

	if len(g.Children["myapp.config.load_config"]) != 0 {
		t.Errorf("load_config children = %v, want empty", g.Children["myapp.config.load_config"])
	}

	if len(g.Parents) != 0 {
		t.Errorf("parents = %v, want empty", g.Parents)
	}
}
