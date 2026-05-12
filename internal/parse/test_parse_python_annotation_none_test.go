//ff:func feature=parse type=util control=sequence
//ff:what test: TestParsePythonAnnotationNone
package parse

import "testing"

func TestParsePythonAnnotationNone(t *testing.T) {
	ann, err := ParsePythonAnnotation("testdata/no_annotation.py")
	if err != nil {
		t.Fatalf("ParsePythonAnnotation failed: %v", err)
	}
	if ann != nil {
		t.Errorf("expected nil annotation, got %+v", ann)
	}
}
