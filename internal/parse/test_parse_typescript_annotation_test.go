//ff:func feature=parse type=util control=sequence
//ff:what test: TestParseTypeScriptAnnotation
package parse

import "testing"

func TestParseTypeScriptAnnotation(t *testing.T) {
	ann, err := ParseTypeScriptAnnotation("testdata/ts_annotated.ts")
	if err != nil {
		t.Fatalf("ParseTypeScriptAnnotation failed: %v", err)
	}
	if ann == nil {
		t.Fatal("annotation is nil")
	}
	if ann.Func["feature"] != "validate" {
		t.Errorf("feature = %q, want %q", ann.Func["feature"], "validate")
	}
	if ann.Func["type"] != "rule" {
		t.Errorf("type = %q, want %q", ann.Func["type"], "rule")
	}
	if ann.Func["control"] != "sequence" {
		t.Errorf("control = %q, want %q", ann.Func["control"], "sequence")
	}
	if ann.What != "TypeScript 유효성 검증 함수" {
		t.Errorf("What = %q, want %q", ann.What, "TypeScript 유효성 검증 함수")
	}
	if ann.Why != "타입 안전성 확보" {
		t.Errorf("Why = %q, want %q", ann.Why, "타입 안전성 확보")
	}
}
