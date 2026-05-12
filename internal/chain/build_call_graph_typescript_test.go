//ff:func feature=chain type=util control=sequence
//ff:what test: TestBuildCallGraphTypeScript
package chain

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func TestBuildCallGraphTypeScript(t *testing.T) {
	files := []model.SourceFile{
		&model.TypeScriptFile{
			Module: "src/service",
			Path:   "src/service.ts",
			Funcs:  []string{"handleRequest"},
			Calls:  []string{"src/utils.validate", "src/utils.logError"},
		},
		&model.TypeScriptFile{
			Module: "src/utils",
			Path:   "src/validate.ts",
			Funcs:  []string{"validate"},
			Calls:  []string{"src/utils.logError"},
		},
		&model.TypeScriptFile{
			Module: "src/utils",
			Path:   "src/log_error.ts",
			Funcs:  []string{"logError"},
			Calls:  nil,
		},
	}

	g := BuildCallGraph(files, "", nil)

	children := g.Children["src/service.handleRequest"]
	if len(children) != 2 {
		t.Fatalf("handleRequest children = %v, want 2 items", children)
	}

	children = g.Children["src/utils.validate"]
	if len(children) != 1 || children[0] != "src/utils.logError" {
		t.Errorf("validate children = %v, want [src/utils.logError]", children)
	}

	if len(g.Children["src/utils.logError"]) != 0 {
		t.Errorf("logError children = %v, want empty", g.Children["src/utils.logError"])
	}

	parents := g.Parents["src/utils.logError"]
	if len(parents) != 2 {
		t.Errorf("logError parents = %v, want 2 items", parents)
	}
}
