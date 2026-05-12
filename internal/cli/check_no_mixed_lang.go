//ff:func feature=validate type=rule control=iteration dimension=1
//ff:what 프로젝트에 타 언어 파일 혼재 시 ERROR 반환 (P1 단일 언어 룰)
package cli

import (
	"fmt"

	"github.com/park-jun-woo/filefunc/internal/walk"
)

// CheckNoMixedLang checks that the project contains only files for the given language.
// Returns an error if files of other languages are found.
func CheckNoMixedLang(root string, lang string, ignorePatterns []string) error {
	foreignExts := foreignExtensions(lang)
	for _, ext := range foreignExts {
		files, err := walk.WalkFiles(root, ext, ignorePatterns)
		if err != nil {
			return fmt.Errorf("walk error for %s: %w", ext, err)
		}
		if len(files) > 0 {
			return fmt.Errorf("mixed language project — found %d %s file(s) in %s project; separate into distinct projects", len(files), ext, lang)
		}
	}
	return nil
}
