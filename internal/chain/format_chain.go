//ff:func feature=chain type=formatter control=iteration dimension=1
//ff:what 체인 탐색 결과를 텍스트로 출력 (--meta, --rate 옵션 지원)
package chain

import (
	"fmt"
	"io"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// FormatChain writes chain results as human-readable text.
// metaFlags: set of requested meta fields (meta, what, why, checked). nil = no meta.
// fileMap: func/method/type name → GoFile mapping. nil = no meta.
// scores: result index → relevance score. nil = no scoring.
// removed: number of results removed by rate filter. 0 = no filtering.
func FormatChain(w io.Writer, start string, results []ChonResult, metaFlags map[string]bool, fileMap map[string]*model.GoFile, scores map[int]float64, removed int) {
	startDisplay := nameFromQualified(start)
	fmt.Fprintf(w, "%s%s\n", startDisplay, formatMeta(start, metaFlags, fileMap))
	startPkg := pkgFromQualified(start)
	for i, r := range results {
		scoreSuffix := ""
		if s, ok := scores[i]; ok {
			scoreSuffix = fmt.Sprintf(" [%.2f]", s)
		}
		display := nameFromQualified(r.Name)
		if pkgFromQualified(r.Name) != startPkg {
			display = r.Name
		}
		fmt.Fprintf(w, "  %d촌 %s: %s%s%s\n", r.Chon, r.Rel, display, formatMeta(r.Name, metaFlags, fileMap), scoreSuffix)
	}
	if len(results) == 0 {
		fmt.Fprintln(w, "  (no relations found)")
	}
	if removed > 0 {
		fmt.Fprintf(w, "  -- %d results filtered by rate (shown %d) --\n", removed, countChon2Plus(results))
	}
}
