//ff:func feature=context type=formatter control=iteration dimension=1
//ff:what context 파이프라인 최종 결과를 출력
package context

import (
	"fmt"
	"io"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// FormatResult writes the final context pipeline results.
func FormatResult(w io.Writer, files []model.SourceFile, scores map[int]float64) {
	fmt.Fprintln(w, "\nResults:")
	for i, sf := range files {
		name := funcName(sf)
		var score float64
		if scores != nil {
			score = scores[i]
		}
		what := ""
		ann := sf.GetAnnotation()
		if ann != nil {
			what = ann.What
		}
		fmt.Fprintf(w, "  %s [%.2f] (what=\"%s\")\n", name, score, what)
	}
	if len(files) == 0 {
		fmt.Fprintln(w, "  (no results)")
	}
}
