//ff:func feature=parse type=util control=sequence
//ff:what test: TestDetectLangPython
package walk

import (
	"os"
	"testing"
)

func TestDetectLangPython(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(dir+"/pyproject.toml", []byte("[project]"), 0644); err != nil {
		t.Fatal(err)
	}
	got := DetectLang(dir)
	if got != "python" {
		t.Errorf("DetectLang = %q, want %q", got, "python")
	}
}
