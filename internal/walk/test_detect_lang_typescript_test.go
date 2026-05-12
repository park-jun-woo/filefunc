//ff:func feature=parse type=util control=sequence
//ff:what test: TestDetectLangTypescript
package walk

import (
	"os"
	"testing"
)

func TestDetectLangTypescript(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(dir+"/tsconfig.json", []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}
	got := DetectLang(dir)
	if got != "typescript" {
		t.Errorf("DetectLang = %q, want %q", got, "typescript")
	}
}
