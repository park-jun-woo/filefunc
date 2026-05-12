//ff:func feature=parse type=util control=sequence
//ff:what test: TestParsePythonAnnotationShebang
package parse

import "testing"

func TestParsePythonAnnotationShebang(t *testing.T) {
	ann, err := ParsePythonAnnotation("testdata/shebang_encoding.py")
	if err != nil {
		t.Fatalf("ParsePythonAnnotation failed: %v", err)
	}
	if ann == nil {
		t.Fatal("annotation is nil")
	}
	if ann.Func["feature"] != "parse" {
		t.Errorf("feature = %q, want %q", ann.Func["feature"], "parse")
	}
	if ann.What != "shebang과 encoding 선언 이후 어노테이션" {
		t.Errorf("What = %q, want %q", ann.What, "shebang과 encoding 선언 이후 어노테이션")
	}
}
