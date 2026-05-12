//ff:func feature=cli type=util control=selection
//ff:what 언어에 따라 어노테이션 접두사를 반환 — Go는 "//ff:", Python은 "# ff:"
package cli

import "github.com/park-jun-woo/filefunc/internal/model"

// AnnotationPrefixForLang returns the annotation prefix string for error messages.
func AnnotationPrefixForLang(sf model.SourceFile) string {
	switch sf.GetLang() {
	case "python":
		return "# ff:"
	default:
		return "//ff:"
	}
}
