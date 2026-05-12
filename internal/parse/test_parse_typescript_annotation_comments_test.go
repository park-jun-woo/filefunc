//ff:func feature=parse type=util control=sequence
//ff:what test: TestParseTypeScriptAnnotationComments
package parse

import "testing"

func TestParseTypeScriptAnnotationComments(t *testing.T) {
	ann, err := ParseTypeScriptAnnotation("testdata/ts_with_comments.ts")
	if err != nil {
		t.Fatalf("ParseTypeScriptAnnotation failed: %v", err)
	}
	if ann == nil {
		t.Fatal("annotation is nil")
	}
	if ann.Func["feature"] != "parse" {
		t.Errorf("feature = %q, want %q", ann.Func["feature"], "parse")
	}
	if ann.Func["type"] != "parser" {
		t.Errorf("type = %q, want %q", ann.Func["type"], "parser")
	}
	if ann.Func["control"] != "iteration" {
		t.Errorf("control = %q, want %q", ann.Func["control"], "iteration")
	}
	if ann.Func["dimension"] != "2" {
		t.Errorf("dimension = %q, want %q", ann.Func["dimension"], "2")
	}
	if ann.What != "주석 이후 어노테이션 파싱 테스트" {
		t.Errorf("What = %q, want %q", ann.What, "주석 이후 어노테이션 파싱 테스트")
	}
}
