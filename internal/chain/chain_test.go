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
