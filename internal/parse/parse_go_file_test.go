package parse

import (
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
