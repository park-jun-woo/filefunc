//ff:func feature=chain type=formatter control=iteration dimension=1
//ff:what 체인 탐색 결과를 텍스트로 출력
package chain

import (
	"fmt"
	"io"
)

// FormatChain writes chain results as human-readable text.
func FormatChain(w io.Writer, start string, results []ChonResult) {
	fmt.Fprintf(w, "%s\n", start)
	for _, r := range results {
		fmt.Fprintf(w, "  %d촌 %s: %s\n", r.Chon, r.Rel, r.Name)
	}
	if len(results) == 0 {
		fmt.Fprintln(w, "  (no relations found)")
	}
}
