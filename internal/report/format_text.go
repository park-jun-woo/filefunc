//ff:func feature=report type=formatter
//ff:what 검증 위반 목록을 텍스트 형식으로 출력
package report

import (
	"fmt"
	"io"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// FormatText writes violations as human-readable text to the given writer.
func FormatText(w io.Writer, violations []model.Violation) {
	if len(violations) == 0 {
		fmt.Fprintln(w, "No violations found.")
		return
	}
	for _, v := range violations {
		fmt.Fprintf(w, "[%s] %s: %s (%s)\n", v.Level, v.Rule, v.Message, v.File)
	}
	fmt.Fprintf(w, "\n%d violation(s) found.\n", len(violations))
}
