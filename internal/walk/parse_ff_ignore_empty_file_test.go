//ff:func feature=parse type=util control=sequence
//ff:what test: TestParseFFIgnore_EmptyFile
package walk

import (
	"os"
	"testing"
)

func TestParseFFIgnore_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/.ffignore"
	os.WriteFile(path, []byte("# comment\n\nvendor/\n*.tmp\n"), 0644)
	patterns := ParseFFIgnore(path)
	if len(patterns) != 2 {
		t.Fatalf("len = %d, want 2", len(patterns))
	}
	if patterns[0] != "vendor/" || patterns[1] != "*.tmp" {
		t.Errorf("patterns = %v, want [vendor/ *.tmp]", patterns)
	}
}
