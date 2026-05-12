//ff:func feature=chain type=util control=iteration dimension=1
//ff:what test: TestBuildCallGraphMixed — Go와 Python 혼합 시 각 언어 내부만 그래프 구축
package chain

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func TestBuildCallGraphMixed(t *testing.T) {
	files := []model.SourceFile{
		&model.GoFile{
			Package: "testdata",
			Path:    "testdata/caller.go",
			Funcs:   []string{"Caller"},
		},
		&model.GoFile{
			Package: "testdata",
			Path:    "testdata/helper_a.go",
			Funcs:   []string{"HelperA"},
		},
		&model.PythonFile{
			Module: "myapp.service",
			Path:   "myapp/service.py",
			Funcs:  []string{"handle"},
			Calls:  []string{"myapp.utils.validate"},
		},
		&model.PythonFile{
			Module: "myapp.utils",
			Path:   "myapp/utils.py",
			Funcs:  []string{"validate"},
			Calls:  nil,
		},
	}

	projFuncs := map[string]string{
		"Caller":  "testdata",
		"HelperA": "testdata",
	}

	g := BuildCallGraph(files, "github.com/nonexistent", projFuncs)

	// Python edges exist
	pyChildren := g.Children["myapp.service.handle"]
	if len(pyChildren) != 1 || pyChildren[0] != "myapp.utils.validate" {
		t.Errorf("Python handle children = %v, want [myapp.utils.validate]", pyChildren)
	}

	// Python validate has no children
	if len(g.Children["myapp.utils.validate"]) != 0 {
		t.Errorf("Python validate children = %v, want empty", g.Children["myapp.utils.validate"])
	}

	// Go Caller's edges are built via ExtractCalls (may or may not find things in testdata)
	// The key check: Go and Python graphs are independent
	// Python funcs do not appear as children of Go funcs
	goChildren := g.Children["testdata.Caller"]
	for _, c := range goChildren {
		if c == "myapp.service.handle" || c == "myapp.utils.validate" {
			t.Errorf("Go Caller should not call Python func %s", c)
		}
	}

	// Python parents don't include Go funcs
	pyParents := g.Parents["myapp.utils.validate"]
	for _, p := range pyParents {
		if p == "testdata.Caller" || p == "testdata.HelperA" {
			t.Errorf("Python validate should not have Go parent %s", p)
		}
	}
}
