//ff:func feature=validate type=rule control=sequence
//ff:what C2: codebook.yaml 원본 텍스트에서 동일 섹션 내 중복 키를 검출
package validate

import (
	"os"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// CheckCodebookDuplicates checks C2: no duplicate keys within the same section.
func CheckCodebookDuplicates(codebookPath string) []model.Violation {
	data, err := os.ReadFile(codebookPath)
	if err != nil {
		return nil
	}
	return FindDuplicateKeys(string(data))
}
