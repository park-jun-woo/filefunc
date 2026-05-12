//ff:func feature=cli type=util control=sequence
//ff:what --lang 플래그 값을 해석하고 비어있으면 DetectLang으로 자동 감지
package cli

import (
	"fmt"

	"github.com/park-jun-woo/filefunc/internal/walk"
)

// ResolveLang returns the language from the flag value, or auto-detects from root.
func ResolveLang(flagValue string, root string) (string, error) {
	lang := flagValue
	if lang == "" {
		lang = walk.DetectLang(root)
	}
	if lang == "" {
		return "", fmt.Errorf("cannot detect language; use --lang go or --lang python")
	}
	return lang, nil
}
