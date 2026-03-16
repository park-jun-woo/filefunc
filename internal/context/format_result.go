//ff:func feature=context type=formatter control=iteration dimension=1
//ff:what context 파이프라인 최종 결과를 출력
package context

import (
	"fmt"
	"io"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// FormatResult writes the final context pipeline results.
func FormatResult(w io.Writer, files []*model.GoFile, scores map[int]float64) {
	fmt.Fprintln(w, "\nResults:")
	for i, gf := range files {
		name := funcName(gf)
		score := scores[i]
		what := ""
		if gf.Annotation != nil {
			what = gf.Annotation.What
		}
		fmt.Fprintf(w, "  %s [%.2f] (what=\"%s\")\n", name, score, what)
	}
	if len(files) == 0 {
		fmt.Fprintln(w, "  (no results)")
	}
}
