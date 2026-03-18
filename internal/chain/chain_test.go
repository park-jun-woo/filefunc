package chain

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func TestNameFromQualified(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"pkg.FuncName", "FuncName"},
		{"a.b.c.Func", "Func"},
		{"NoPackage", "NoPackage"},
		{"", ""},
	}
	for _, tt := range tests {
		got := NameFromQualified(tt.input)
		if got != tt.want {
			t.Errorf("NameFromQualified(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestPkgFromQualified(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"pkg.FuncName", "pkg"},
		{"a.b.c.Func", "a.b.c"},
		{"NoPackage", ""},
		{"", ""},
	}
	for _, tt := range tests {
		got := PkgFromQualified(tt.input)
		if got != tt.want {
			t.Errorf("PkgFromQualified(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestQualifiedName(t *testing.T) {
	got := qualifiedName("validate", "RuleF1")
	if got != "validate.RuleF1" {
		t.Errorf("qualifiedName = %q, want %q", got, "validate.RuleF1")
	}
}

func TestFuncName(t *testing.T) {
	tests := []struct {
		name string
		gf   *model.GoFile
		want string
	}{
		{"func", &model.GoFile{Funcs: []string{"Foo"}}, "Foo"},
		{"method", &model.GoFile{Methods: []string{"Bar"}}, "Bar"},
		{"func over method", &model.GoFile{Funcs: []string{"Foo"}, Methods: []string{"Bar"}}, "Foo"},
		{"empty", &model.GoFile{}, ""},
	}
	for _, tt := range tests {
		got := funcName(tt.gf)
		if got != tt.want {
			t.Errorf("funcName(%s) = %q, want %q", tt.name, got, tt.want)
		}
	}
}

func TestAddUnique(t *testing.T) {
	seen := map[string]bool{"a": true}
	var result []string
	AddUnique([]string{"a", "b", "c", "b"}, seen, &result)
	if len(result) != 2 {
		t.Fatalf("len = %d, want 2", len(result))
	}
	if result[0] != "b" || result[1] != "c" {
		t.Errorf("result = %v, want [b c]", result)
	}
}

func TestCollectChon(t *testing.T) {
	seen := map[string]bool{"a": true}
	results := CollectChon([]string{"a", "b", "c"}, 2, "child", seen)
	if len(results) != 2 {
		t.Fatalf("len = %d, want 2", len(results))
	}
	if results[0].Name != "b" || results[0].Chon != 2 || results[0].Rel != "child" {
		t.Errorf("results[0] = %+v", results[0])
	}
}

func TestCountChon2Plus(t *testing.T) {
	results := []ChonResult{
		{"a", 1, "child"},
		{"b", 2, "grandchild"},
		{"c", 3, "grandchild"},
	}
	got := countChon2Plus(results)
	if got != 2 {
		t.Errorf("countChon2Plus = %d, want 2", got)
	}
}

func TestHasFeature(t *testing.T) {
	gf := &model.GoFile{Annotation: &model.Annotation{
		Func: map[string]string{"feature": "validate"},
	}}
	if !hasFeature(gf, "validate") {
		t.Error("expected true for feature=validate")
	}
	if hasFeature(gf, "chain") {
		t.Error("expected false for feature=chain")
	}
	if hasFeature(&model.GoFile{}, "validate") {
		t.Error("expected false for nil annotation")
	}
}

func TestFilterByFeature(t *testing.T) {
	files := []*model.GoFile{
		{Package: "v", Funcs: []string{"RuleF1"}, Annotation: &model.Annotation{Func: map[string]string{"feature": "validate"}}},
		{Package: "c", Funcs: []string{"Build"}, Annotation: &model.Annotation{Func: map[string]string{"feature": "chain"}}},
	}
	got := FilterByFeature(files, "validate")
	if len(got) != 1 || got[0] != "v.RuleF1" {
		t.Errorf("FilterByFeature = %v, want [v.RuleF1]", got)
	}
}

func TestFilterByPackage(t *testing.T) {
	results := []ChonResult{
		{"validate.RuleF1", 1, "child"},
		{"chain.Build", 2, "grandchild"},
	}
	got := FilterByPackage(results, "validate")
	if len(got) != 1 || got[0].Name != "validate.RuleF1" {
		t.Errorf("FilterByPackage = %v", got)
	}
}

// --- Graph integration tests ---

func buildTestGraph() *CallGraph {
	// Caller -> HelperA, HelperB
	// HelperA -> Leaf
	g := &CallGraph{
		Children: map[string][]string{
			"testdata.Caller":  {"testdata.HelperA", "testdata.HelperB"},
			"testdata.HelperA": {"testdata.Leaf"},
		},
		Parents: map[string][]string{
			"testdata.HelperA": {"testdata.Caller"},
			"testdata.HelperB": {"testdata.Caller"},
			"testdata.Leaf":    {"testdata.HelperA"},
		},
	}
	return g
}

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

func TestTraverseChon_Chon1(t *testing.T) {
	g := buildTestGraph()
	results := TraverseChon(g, "testdata.HelperA", 1)

	hasChild := false
	hasParent := false
	for _, r := range results {
		if r.Name == "testdata.Leaf" && r.Chon == 1 && r.Rel == "calls" {
			hasChild = true
		}
		if r.Name == "testdata.Caller" && r.Chon == 1 && r.Rel == "called-by" {
			hasParent = true
		}
	}
	if !hasChild {
		t.Error("missing chon=1 child: testdata.Leaf")
	}
	if !hasParent {
		t.Error("missing chon=1 parent: testdata.Caller")
	}
}

func TestTraverseChon_Chon2(t *testing.T) {
	g := buildTestGraph()
	results := TraverseChon(g, "testdata.HelperA", 2)

	hasSibling := false
	for _, r := range results {
		if r.Name == "testdata.HelperB" && r.Chon == 2 && r.Rel == "co-called" {
			hasSibling = true
		}
	}
	if !hasSibling {
		t.Error("missing chon=2 sibling: testdata.HelperB")
	}
}

func TestTraverseDepth_Children(t *testing.T) {
	g := buildTestGraph()
	results := TraverseDepth(g, "testdata.Caller", "calls", 3)

	names := make(map[string]bool)
	for _, r := range results {
		names[r.Name] = true
	}
	if !names["testdata.HelperA"] || !names["testdata.HelperB"] || !names["testdata.Leaf"] {
		t.Errorf("TraverseDepth children = %v, want HelperA, HelperB, Leaf", results)
	}
}

func TestTraverseDepth_Parents(t *testing.T) {
	g := buildTestGraph()
	results := TraverseDepth(g, "testdata.Leaf", "called-by", 3)

	names := make(map[string]bool)
	for _, r := range results {
		names[r.Name] = true
	}
	if !names["testdata.HelperA"] || !names["testdata.Caller"] {
		t.Errorf("TraverseDepth parents = %v, want HelperA, Caller", results)
	}
}

func TestFindSiblings(t *testing.T) {
	g := buildTestGraph()
	siblings := FindSiblings(g, "testdata.HelperA")
	if len(siblings) != 1 || siblings[0] != "testdata.HelperB" {
		t.Errorf("FindSiblings = %v, want [testdata.HelperB]", siblings)
	}
}

func TestExpandThrough(t *testing.T) {
	g := buildTestGraph()
	// From Caller's children, expand their children
	result := ExpandThrough(g.Children["testdata.Caller"], func(c string) []string {
		return g.Children[c]
	})
	if len(result) != 1 || result[0] != "testdata.Leaf" {
		t.Errorf("ExpandThrough = %v, want [testdata.Leaf]", result)
	}
}
