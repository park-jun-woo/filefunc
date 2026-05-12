//ff:func feature=chain type=util control=sequence
//ff:what test: TestBuildCallGraphPython
package chain

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func TestBuildCallGraphPython(t *testing.T) {
	files := []model.SourceFile{
		&model.PythonFile{
			Module: "myapp.service",
			Path:   "myapp/service.py",
			Funcs:  []string{"handle_request"},
			Calls:  []string{"myapp.utils.validate", "myapp.utils.log_error"},
		},
		&model.PythonFile{
			Module: "myapp.utils",
			Path:   "myapp/utils.py",
			Funcs:  []string{"validate"},
			Calls:  []string{"myapp.utils.log_error"},
		},
		&model.PythonFile{
			Module: "myapp.utils",
			Path:   "myapp/log_error.py",
			Funcs:  []string{"log_error"},
			Calls:  nil,
		},
	}

	g := BuildCallGraph(files, "", nil)

	// service.handle_request → utils.validate, utils.log_error
	children := g.Children["myapp.service.handle_request"]
	if len(children) != 2 {
		t.Fatalf("handle_request children = %v, want 2 items", children)
	}

	// utils.validate → utils.log_error
	children = g.Children["myapp.utils.validate"]
	if len(children) != 1 || children[0] != "myapp.utils.log_error" {
		t.Errorf("validate children = %v, want [myapp.utils.log_error]", children)
	}

	// log_error has no children
	if len(g.Children["myapp.utils.log_error"]) != 0 {
		t.Errorf("log_error children = %v, want empty", g.Children["myapp.utils.log_error"])
	}

	// Parents: log_error's parents are handle_request and validate
	parents := g.Parents["myapp.utils.log_error"]
	if len(parents) != 2 {
		t.Errorf("log_error parents = %v, want 2 items", parents)
	}
}
