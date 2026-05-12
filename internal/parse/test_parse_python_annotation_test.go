//ff:func feature=parse type=util control=sequence
//ff:what test: TestParsePythonAnnotation
package parse

import "testing"

func TestParsePythonAnnotation(t *testing.T) {
	ann, err := ParsePythonAnnotation("testdata/annotated.py")
	if err != nil {
		t.Fatalf("ParsePythonAnnotation failed: %v", err)
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
	if ann.What != "유효성 검증 룰 적용" {
		t.Errorf("What = %q, want %q", ann.What, "유효성 검증 룰 적용")
	}
}
