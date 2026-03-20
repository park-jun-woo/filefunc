//ff:func feature=parse type=util control=sequence
//ff:what test: TestParseGoFileAnnotation
package parse

import (
	"testing"
)

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
