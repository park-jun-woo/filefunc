//ff:func feature=chain type=util control=sequence
//ff:what test: TestBuildCallGraph
package chain

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func TestBuildCallGraph(t *testing.T) {
	files := []*model.GoFile{
		{Package: "testdata", Path: "testdata/caller.go", Funcs: []string{"Caller"}},
		{Package: "testdata", Path: "testdata/helper_a.go", Funcs: []string{"HelperA"}},
		{Package: "testdata", Path: "testdata/helper_b.go", Funcs: []string{"HelperB"}},
		{Package: "testdata", Path: "testdata/leaf.go", Funcs: []string{"Leaf"}},
	}
	projFuncs := map[string]string{
		"Caller":  "testdata",
		"HelperA": "testdata",
		"HelperB": "testdata",
		"Leaf":    "testdata",
	}
	g := BuildCallGraph(files, "github.com/nonexistent", projFuncs)

	// Caller should call HelperA and HelperB
	children := g.Children["testdata.Caller"]
	if len(children) != 2 {
		t.Fatalf("Caller children = %v, want 2 items", children)
	}

	// HelperA should call Leaf
	children = g.Children["testdata.HelperA"]
	if len(children) != 1 || children[0] != "testdata.Leaf" {
		t.Errorf("HelperA children = %v, want [testdata.Leaf]", children)
	}

	// Leaf has no children
	if len(g.Children["testdata.Leaf"]) != 0 {
		t.Errorf("Leaf children = %v, want empty", g.Children["testdata.Leaf"])
	}

	// Parents: HelperA's parent is Caller
	parents := g.Parents["testdata.HelperA"]
	if len(parents) != 1 || parents[0] != "testdata.Caller" {
		t.Errorf("HelperA parents = %v, want [testdata.Caller]", parents)
	}
}
