//ff:func feature=context type=formatter control=iteration dimension=1
//ff:what context 파이프라인 최종 결과를 출력
package context

import (
	"fmt"
	"io"

	"github.com/park-jun-woo/filefunc/internal/chain"
	"github.com/park-jun-woo/filefunc/internal/model"
)

// FormatResult writes the final context pipeline results.
func FormatResult(w io.Writer, results []chain.ChonResult, scores map[int]float64, fileMap map[string]*model.GoFile) {
	fmt.Fprintln(w, "\nResults:")
	for i, r := range results {
		score := scores[i]
		what := ""
		if gf := fileMap[r.Name]; gf != nil && gf.Annotation != nil {
			what = gf.Annotation.What
		}
		fmt.Fprintf(w, "  %s [%.2f] %d촌 (what=\"%s\")\n", r.Name, score, r.Chon, what)
	}
	if len(results) == 0 {
		fmt.Fprintln(w, "  (no results)")
	}
}
