package parse

import (
	"os"
	"testing"
)

func TestParseGoFile(t *testing.T) {
	gf, err := ParseGoFile("testdata/sample.go")
	if err != nil {
		t.Fatalf("ParseGoFile failed: %v", err)
	}

	if gf.Package != "testdata" {
		t.Errorf("Package = %q, want %q", gf.Package, "testdata")
	}

	if len(gf.Funcs) != 1 || gf.Funcs[0] != "CheckSample" {
		t.Errorf("Funcs = %v, want [CheckSample]", gf.Funcs)
	}

	if len(gf.Types) != 1 || gf.Types[0] != "SampleParam" {
		t.Errorf("Types = %v, want [SampleParam]", gf.Types)
	}

	if gf.HasInit {
		t.Error("HasInit = true, want false")
	}

	if gf.MaxDepth != 2 {
		t.Errorf("MaxDepth = %d, want 2", gf.MaxDepth)
	}

	if gf.IsTest {
		t.Error("IsTest = true, want false")
	}
}

func TestParseGoFileAnnotation(t *testing.T) {
	gf, err := ParseGoFile("testdata/sample.go")
	if err != nil {
		t.Fatalf("ParseGoFile failed: %v", err)
	}

	if gf.Annotation == nil {
		t.Fatal("Annotation is nil")
	}

	if gf.Annotation.Func["feature"] != "validate" {
		t.Errorf("feature = %q, want %q", gf.Annotation.Func["feature"], "validate")
	}

	if gf.Annotation.What != "파일당 func 개수를 검증한다" {
		t.Errorf("What = %q, want %q", gf.Annotation.What, "파일당 func 개수를 검증한다")
	}

	if gf.Annotation.Why != "제1시민은 AI 에이전트" {
		t.Errorf("Why = %q, want %q", gf.Annotation.Why, "제1시민은 AI 에이전트")
	}

	// calls/uses removed — computed on demand via AST
}

func TestCallName_SamePackage(t *testing.T) {
	projFuncs := map[string]string{"HelperA": "testdata", "Leaf": "testdata"}
	projImports := map[string]string{}

	calls, err := ExtractCalls("testdata/caller.go", "github.com/nonexistent", projFuncs, "testdata")
	if err != nil {
		t.Fatalf("ExtractCalls failed: %v", err)
	}

	found := make(map[string]bool)
	for _, c := range calls {
		found[c] = true
	}
	if !found["testdata.HelperA"] {
		t.Errorf("missing call testdata.HelperA, got %v", calls)
	}
	_ = projImports
}

func TestCallName_ExternalIgnored(t *testing.T) {
	projFuncs := map[string]string{}
	calls, err := ExtractCalls("testdata/caller.go", "github.com/nonexistent", projFuncs, "testdata")
	if err != nil {
		t.Fatalf("ExtractCalls failed: %v", err)
	}
	if len(calls) != 0 {
		t.Errorf("expected no calls for unknown funcs, got %v", calls)
	}
}

func TestDetectControl_Sequence(t *testing.T) {
	got := DetectControl("testdata/detect_sequence.go")
	if got != "sequence" {
		t.Errorf("DetectControl(sequence) = %q, want %q", got, "sequence")
	}
}

func TestDetectControl_Selection(t *testing.T) {
	got := DetectControl("testdata/detect_selection.go")
	if got != "selection" {
		t.Errorf("DetectControl(selection) = %q, want %q", got, "selection")
	}
}

func TestDetectControl_Iteration(t *testing.T) {
	got := DetectControl("testdata/detect_iteration.go")
	if got != "iteration" {
		t.Errorf("DetectControl(iteration) = %q, want %q", got, "iteration")
	}
}

func TestCalcMaxDepth_Nested(t *testing.T) {
	gf, err := ParseGoFile("testdata/depth_nested.go")
	if err != nil {
		t.Fatalf("ParseGoFile failed: %v", err)
	}
	// for { if { if {} } } = depth 3
	if gf.MaxDepth != 3 {
		t.Errorf("MaxDepth = %d, want 3", gf.MaxDepth)
	}
}

func TestReadModulePath_Missing(t *testing.T) {
	_, err := ReadModulePath("/nonexistent/go.mod")
	if err == nil {
		t.Error("expected error for nonexistent go.mod")
	}
}

func TestReadModulePath_NoModule(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/go.mod"
	if err := writeTestFile(path, "go 1.22\n"); err != nil {
		t.Fatal(err)
	}
	_, err := ReadModulePath(path)
	if err == nil {
		t.Error("expected error for go.mod without module directive")
	}
}

func TestReadModulePath_Valid(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/go.mod"
	if err := writeTestFile(path, "module github.com/test/proj\n\ngo 1.22\n"); err != nil {
		t.Fatal(err)
	}
	mod, err := ReadModulePath(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mod != "github.com/test/proj" {
		t.Errorf("module = %q, want %q", mod, "github.com/test/proj")
	}
}

func writeTestFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}
