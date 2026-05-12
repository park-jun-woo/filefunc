//ff:func feature=parse type=util control=sequence
//ff:what test: TestParseTypeScriptAnnotationNone
package parse

import "testing"

func TestParseTypeScriptAnnotationNone(t *testing.T) {
	ann, err := ParseTypeScriptAnnotation("testdata/ts_no_annotation.ts")
	if err != nil {
		t.Fatalf("ParseTypeScriptAnnotation failed: %v", err)
	}
	if ann != nil {
		t.Errorf("expected nil annotation, got %+v", ann)
	}
}
