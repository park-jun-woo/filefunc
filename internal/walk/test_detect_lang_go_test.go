//ff:func feature=parse type=util control=sequence
//ff:what test: TestDetectLangGo
package walk

import (
	"os"
	"testing"
)

func TestDetectLangGo(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(dir+"/go.mod", []byte("module test"), 0644); err != nil {
		t.Fatal(err)
	}
	got := DetectLang(dir)
	if got != "go" {
		t.Errorf("DetectLang = %q, want %q", got, "go")
	}
}
