//ff:func feature=report type=formatter
//ff:what 검증 위반 목록을 JSON 형식으로 출력
package report

import (
	"encoding/json"
	"io"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// FormatJSON writes violations as JSON to the given writer.
func FormatJSON(w io.Writer, violations []model.Violation) error {
	if violations == nil {
		violations = []model.Violation{}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(violations)
}
