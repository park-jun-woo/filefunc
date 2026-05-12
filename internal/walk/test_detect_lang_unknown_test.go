//ff:func feature=parse type=util control=sequence
//ff:what test: TestDetectLangUnknown
package walk

import "testing"

func TestDetectLangUnknown(t *testing.T) {
	dir := t.TempDir()
	got := DetectLang(dir)
	if got != "" {
		t.Errorf("DetectLang = %q, want %q", got, "")
	}
}
