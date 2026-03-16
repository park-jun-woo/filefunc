//ff:func feature=chain type=formatter control=iteration dimension=1
//ff:what 체인 탐색 결과를 텍스트로 출력 (--meta 옵션 지원)
package chain

import (
	"fmt"
	"io"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// FormatChain writes chain results as human-readable text.
// metaFlags: set of requested meta fields (meta, what, why, checked). nil = no meta.
// fileMap: func/method/type name → GoFile mapping. nil = no meta.
func FormatChain(w io.Writer, start string, results []ChonResult, metaFlags map[string]bool, fileMap map[string]*model.GoFile) {
	fmt.Fprintf(w, "%s%s\n", start, formatMeta(start, metaFlags, fileMap))
	for _, r := range results {
		fmt.Fprintf(w, "  %d촌 %s: %s%s\n", r.Chon, r.Rel, r.Name, formatMeta(r.Name, metaFlags, fileMap))
	}
	if len(results) == 0 {
		fmt.Fprintln(w, "  (no relations found)")
	}
}
