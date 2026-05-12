//ff:func feature=chain type=formatter control=sequence
//ff:what 함수명에 해당하는 어노테이션 메타데이터를 한 줄 괄호 포맷으로 생성
package chain

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// formatMeta returns a parenthesized metadata string for the given name.
// Returns empty string if no meta requested or name not found.
func formatMeta(name string, metaFlags map[string]bool, fileMap map[string]model.SourceFile) string {
	if len(metaFlags) == 0 || fileMap == nil {
		return ""
	}
	sf := fileMap[name]
	if sf == nil {
		return ""
	}
	ann := sf.GetAnnotation()
	if ann == nil {
		return ""
	}
	var parts []string
	if metaFlags["meta"] {
		parts = append(parts, metaPairs(ann)...)
	}
	if metaFlags["what"] && ann.What != "" {
		parts = append(parts, fmt.Sprintf("what=%q", ann.What))
	}
	if metaFlags["why"] && ann.Why != "" {
		parts = append(parts, fmt.Sprintf("why=%q", ann.Why))
	}
	if metaFlags["checked"] && len(ann.Checked) > 0 {
		parts = append(parts, fmt.Sprintf("checked=%q", checkedString(ann.Checked)))
	}
	if len(parts) == 0 {
		return ""
	}
	return " (" + strings.Join(parts, " ") + ")"
}
